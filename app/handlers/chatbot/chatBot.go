package chatbot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/app/handlers/script"
	"github.com/stealcash/AgentFlow/db"
	"os"
	"path/filepath"
	"strconv"
)

func CreateChatbot(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotName := c.PostForm("chatbot_name")
	defaultMessage := c.PostForm("default_message")

	// Reuse exported GenerateAPIKey from script package
	publicAPIKey, err := script.GenerateAPIKey()
	if err != nil {
		exceptions.Internal("Failed to generate API key")
		return
	}

	var logoPath string
	file, err := c.FormFile("logo")
	if err == nil {
		uploadDir := fmt.Sprintf("uploads/user_%d/chatbots", userID)
		os.MkdirAll(uploadDir, os.ModePerm)

		filename := fmt.Sprintf("chatbot_logo_%d_%s", userID, filepath.Base(file.Filename))
		logoPath = filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, logoPath); err != nil {
			exceptions.Internal("Logo upload failed")
			return
		}
	}

	query := `
		INSERT INTO chatbots (user_id, chatbot_name, logo_path, default_message, public_api_key)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var chatbotID int
	err = db.DB.QueryRow(query, userID, chatbotName, logoPath, defaultMessage, publicAPIKey).Scan(&chatbotID)
	if err != nil {
		exceptions.Internal("Failed to create chatbot")
		return
	}

	responses.Success(c, "Chatbot created successfully", gin.H{
		"id":             chatbotID,
		"logo":           logoPath,
		"public_api_key": publicAPIKey,
	})
}

func ListChatbots(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := db.DB.Query(`
		SELECT id, chatbot_name, logo_path, default_message
		FROM chatbots
		WHERE user_id = $1
	`, userID)
	if err != nil {
		exceptions.Internal("Failed to fetch chatbots")
	}
	defer rows.Close()

	var chatbots []gin.H
	for rows.Next() {
		var id int
		var name, logo, message string
		rows.Scan(&id, &name, &logo, &message)

		chatbots = append(chatbots, gin.H{
			"id":              id,
			"chatbot_name":    name,
			"logo_path":       logo,
			"default_message": message,
		})
	}

	responses.Success(c, "Chatbots fetched successfully", gin.H{
		"chatbots": chatbots,
	})
}

func DeleteChatbot(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid chatbot ID")
	}

	_, err = db.DB.Exec(`DELETE FROM chatbots WHERE id = $1 AND user_id = $2`, chatbotID, userID)
	if err != nil {
		exceptions.Internal("Failed to delete chatbot")
	}

	responses.Success(c, "Chatbot deleted successfully", nil)
}
