package chatbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/app/globals"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/db"
)

type GeneralQuestionInput struct {
	Question string `json:"question_text" binding:"required"`
	Answer   string `json:"answer_text"`
}

func AddGeneralQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")

	// Get chatbot ID
	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid chatbot ID")
		return
	}

	// Bind input JSON
	var input GeneralQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptions.BadRequest(err.Error())
		return
	}

	// Confirm chatbot belongs to user
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2
		)`, chatbotID, userID).Scan(&exists)
	if err != nil || !exists {
		exceptions.Unauthorized("Invalid chatbot or unauthorized")
		return
	}

	// Insert into PostgreSQL and get ID
	var id int
	err = db.DB.QueryRow(`
		INSERT INTO general_questions (chatbot_id, question_text, answer_text)
		VALUES ($1, $2, $3)
		RETURNING id`,
		chatbotID, input.Question, input.Answer).Scan(&id)
	if err != nil {
		exceptions.Internal("Failed to add general question")
		return
	}

	// If Elastic is enabled, index too
	if globals.Config.ElasticDatabase.RequiredElasticConnection {
		indexGeneralQuestionES(id, chatbotID, input.Question, input.Answer)
	}

	responses.Success(c, "General question added", gin.H{
		"id": id,
	})
}

func indexGeneralQuestionES(id int, chatbotID int, question string, answer string) {
	doc := map[string]interface{}{
		"id":            id,
		"chatbot_id":    chatbotID,
		"question_text": question,
		"answer_text":   answer,
	}

	body, _ := json.Marshal(doc)

	// Use compound ID: chatbotID-generalQuestionID
	docID := fmt.Sprintf("%d-%d", chatbotID, id)

	res, err := db.EsDB.Index(
		"general_questions",
		bytes.NewReader(body),
		db.EsDB.Index.WithDocumentID(docID),
		db.EsDB.Index.WithContext(context.Background()),
	)
	if err != nil || res.IsError() {
		fmt.Printf("Elasticsearch indexing failed: %v\n", err)
	}
	defer res.Body.Close()
}

func GetGeneralQuestions(c *gin.Context) {
	userID := c.GetInt("user_id")

	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid chatbot ID")
		return
	}

	// ✅ Confirm ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2
		)`, chatbotID, userID).Scan(&exists)
	if err != nil || !exists {
		exceptions.Unauthorized("Unauthorized")
		return
	}

	// ✅ Always do PG query FIRST
	questions := getGeneralQuestionsPGOnly(chatbotID)

	// ✅ If Elastic is enabled, run ES check too (result ignored)
	if globals.Config.ElasticDatabase.RequiredElasticConnection {
		go checkGeneralQuestionsES(chatbotID) // optional: run in goroutine so it doesn't block
	}

	responses.Success(c, "Questions fetched", gin.H{
		"questions": questions,
	})
}

func getGeneralQuestionsPGOnly(chatbotID int) []gin.H {
	rows, err := db.DB.Query(`
		SELECT id, question_text, answer_text 
		FROM general_questions
		WHERE chatbot_id = $1
	`, chatbotID)

	if err != nil {
		// You can decide: panic or return empty
		fmt.Printf("PG fetch failed: %v\n", err)
		return []gin.H{}
	}
	defer rows.Close()

	var questions []gin.H
	for rows.Next() {
		var id int
		var questionText, answerText string
		rows.Scan(&id, &questionText, &answerText)

		questions = append(questions, gin.H{
			"id":            id,
			"question_text": questionText,
			"answer_text":   answerText,
		})
	}

	return questions
}

func checkGeneralQuestionsES(chatbotID int) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"chatbot_id": chatbotID,
			},
		},
	}

	body, _ := json.Marshal(query)
	res, err := db.EsDB.Search(
		db.EsDB.Search.WithIndex("general_questions"),
		db.EsDB.Search.WithBody(bytes.NewReader(body)),
		db.EsDB.Search.WithContext(context.Background()),
	)
	if err != nil || res.IsError() {
		fmt.Printf("Elasticsearch search failed for chatbot_id=%d: %v\n", chatbotID, err)
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Elasticsearch decode failed: %v\n", err)
		return
	}

	// Optional: log how many hits found
	total := r["hits"].(map[string]interface{})["total"]
	fmt.Printf("Elasticsearch check: %v hits for chatbot_id=%d\n", total, chatbotID)
}

func DeleteGeneralQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")

	// Validate chatbot ID
	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid chatbot ID")
		return
	}

	// Validate question ID
	questionIDStr := c.Param("question_id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid question ID")
		return
	}

	// ✅ Confirm chatbot belongs to user
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2
		)`, chatbotID, userID).Scan(&exists)
	if err != nil || !exists {
		exceptions.Unauthorized("Unauthorized")
		return
	}

	// ✅ Delete from PostgreSQL first
	_, err = db.DB.Exec(`
		DELETE FROM general_questions
		WHERE id = $1 AND chatbot_id = $2`,
		questionID, chatbotID)
	if err != nil {
		exceptions.Internal("Failed to delete question from PostgreSQL")
		return
	}

	// ✅ If Elastic enabled, also delete there (ignore not found)
	if globals.Config.ElasticDatabase.RequiredElasticConnection {
		deleteGeneralQuestionES(chatbotID, questionID)
	}

	responses.Success(c, "Question deleted", nil)
}

func deleteGeneralQuestionES(chatbotID int, questionID int) {
	// Use same compound ID you used in indexing: chatbotID-questionID
	docID := fmt.Sprintf("%d-%d", chatbotID, questionID)

	res, err := db.EsDB.Delete(
		"general_questions",
		docID,
		db.EsDB.Delete.WithContext(context.Background()),
	)
	if err != nil {
		fmt.Printf("Elasticsearch delete error: %v\n", err)
		return
	}
	defer res.Body.Close()

	// Ignore not found: ES returns 404 with {"result": "not_found"}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err == nil {
		if result, ok := r["result"].(string); ok && result == "not_found" {
			fmt.Printf("Elasticsearch: document %s not found, ignoring\n", docID)
		}
	}

	if res.IsError() && res.StatusCode != 404 {
		fmt.Printf("Elasticsearch delete failed: %v\n", res.Status())
	}
}
