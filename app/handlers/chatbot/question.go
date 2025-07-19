package chatbot

import (
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionInput struct {
	CategoryID int    `json:"category_id" binding:"required"`
	Question   string `json:"question_text" binding:"required"`
	Answer     string `json:"answer_text" binding:"required"`
}

func AddQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")

	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	var input QuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM categories
			JOIN chatbots ON categories.chatbot_id = chatbots.id
			WHERE categories.id = $1 AND categories.chatbot_id = $2 AND chatbots.user_id = $3
		)
	`, input.CategoryID, chatbotID, userID).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid category or chatbot"})
		return
	}

	_, err = db.DB.Exec(`
		INSERT INTO questions (category_id, question_text, answer_text)
		VALUES ($1, $2, $3)
	`, input.CategoryID, input.Question, input.Answer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add question"})
		return
	}
	responses.Success(c, "questions", gin.H{"message": "Question added"})

}

// GetQuestionsByCategory retrieves all questions for a category securely
func GetQuestionsByCategory(c *gin.Context) {
	userID := c.GetInt("user_id")

	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatbot ID"})
		return
	}

	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM categories
			JOIN chatbots ON categories.chatbot_id = chatbots.id
			WHERE categories.id = $1 AND categories.chatbot_id = $2 AND chatbots.user_id = $3
		)
	`, categoryID, chatbotID, userID).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, question_text, answer_text
		FROM questions WHERE category_id = $1
	`, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}
	defer rows.Close()

	var questions []gin.H
	for rows.Next() {
		var id int
		var questionText, answerText string
		if err := rows.Scan(&id, &questionText, &answerText); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan failed"})
			return
		}
		questions = append(questions, gin.H{
			"id":            id,
			"question_text": questionText,
			"answer_text":   answerText,
		})
	}

	if questions == nil {
		questions = []gin.H{}
	}

	responses.Success(c, "questions", gin.H{"questions": questions})
}

// DeleteQuestion deletes a question after verifying user owns the chatbot
func DeleteQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")

	questionIDStr := c.Param("id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM questions
			JOIN categories ON questions.category_id = categories.id
			JOIN chatbots ON categories.chatbot_id = chatbots.id
			WHERE questions.id = $1 AND chatbots.user_id = $2
		)
	`, questionID, userID).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or not found"})
		return
	}

	_, err = db.DB.Exec(`DELETE FROM questions WHERE id = $1`, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	responses.Success(c, "questions", nil)

}
