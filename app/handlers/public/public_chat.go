package public

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/app/globals"
	"github.com/stealcash/AgentFlow/app/handlers/analytics"
	"github.com/stealcash/AgentFlow/app/logger"
	"github.com/stealcash/AgentFlow/db"
	"net/http"
	"strings"
)

type ChatProfile struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	ChatbotName    string `json:"chatbot_name"`
	DefaultMessage string `json:"default_message"`
	LogoPath       string `json:"logo_path"`
}

func PublicChatQuery(c *gin.Context) {
	type QueryRequest struct {
		APIKey string `json:"api_key"`
		Query  string `json:"query"`
		Domain string `json:"domain"`
	}

	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.FailedApiResponse(c, http.StatusBadRequest, err)
		return
	}

	// 1 Validated chatbot and user
	var chatbotID, userID int
	err := db.DB.QueryRow(`
		SELECT id, user_id
		FROM chatbots
		WHERE public_api_key = $1`,
		req.APIKey).Scan(&chatbotID, &userID)

	if err != nil {
		responses.FailedApiResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 2 Check general questions config (TOML)
	for _, block := range globals.GeneralQuestions.GeneralQuestions {
		for _, q := range block.Questions {
			score := calculateSimilarity(req.Query, q)
			if score >= 0.7 {
				responses.SuccessPublicResponse(c, gin.H{"answer": block.Answer})
				return
			}
		}
	}

	var matchedQID int
	var finalAnswer, source string

	rows, err := db.DB.Query(`
		SELECT id, question_text, answer_text
		FROM general_questions
		WHERE chatbot_id = $1`, chatbotID)

	if err != nil {
		responses.FailedApiResponse(c, http.StatusInternalServerError, "Failed to fetch general questions")
		return
	}
	defer rows.Close()

	var bestMatch struct {
		ID       int
		Question string
		Answer   string
		Score    float64
	}

	for rows.Next() {
		var id int
		var question, answer string
		if err := rows.Scan(&id, &question, &answer); err != nil {
			continue
		}

		score := calculateSimilarity(req.Query, question)
		if score > bestMatch.Score && score >= 0.5 {
			bestMatch.ID = id
			bestMatch.Question = question
			bestMatch.Answer = answer
			bestMatch.Score = score
		}
	}

	if bestMatch.Score >= 0.5 {
		finalAnswer = bestMatch.Answer
		source = "general_question"
		matchedQID = bestMatch.ID
	} else {
		aiAnswer, err := callChatGPT(req.Query)
		if err != nil {
			finalAnswer = "Sorry, I couldn't find an answer."
		} else {
			finalAnswer = aiAnswer
		}
		source = "chatgpt"
	}

	err = analytics.LogAnalytics(chatbotID, req.Domain, nil, &matchedQID, req.Query, source)
	if err != nil {
		logger.Error("analytics failed %v", err)
	}

	responses.SuccessPublicResponse(c, gin.H{"answer": finalAnswer})
}

func callChatGPT(query string) (string, error) {
	apiKey := globals.Config.ChatGPTModels[0].ChatGPTAPIKey
	if apiKey == "" {
		return "Sorry, no AI integration available.", nil
	}

	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": query},
		},
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "Sorry, I couldn't find an answer.", nil
	}
	defer resp.Body.Close()

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "Sorry, I couldn't parse the answer.", nil
	}

	if len(response.Choices) > 0 {
		return strings.TrimSpace(response.Choices[0].Message.Content), nil
	}

	return "Sorry, I don't have an answer.", nil
}

func calculateSimilarity(s1, s2 string) float64 {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	wordSet := make(map[string]struct{})
	for _, w := range words2 {
		wordSet[w] = struct{}{}
	}

	common := 0
	for _, w := range words1 {
		if _, exists := wordSet[w]; exists {
			common++
		}
	}

	total := len(words1) + len(words2) - common
	if total == 0 {
		return 1.0
	}
	return float64(common) / float64(total)
}

func GetPublicChatbotProfile(c *gin.Context) {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		responses.FailedApiResponse(c, http.StatusBadRequest, "API key required")
		return
	}

	var profile ChatProfile

	err := db.DB.QueryRow(`
		SELECT id, user_id, chatbot_name, default_message, logo_path
		FROM chatbots
		WHERE public_api_key = $1`,
		apiKey).Scan(&profile.ID, &profile.UserID, &profile.ChatbotName, &profile.DefaultMessage, &profile.LogoPath)

	if err != nil {
		responses.FailedApiResponse(c, http.StatusNotFound, "Chatbot profile not found")
		return
	}

	responses.SuccessPublicResponse(c, profile)
}
