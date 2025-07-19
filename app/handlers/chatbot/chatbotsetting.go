package chatbot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/db"
	"os"
	"path/filepath"
)

func UploadChatbotSettings(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotID := c.Param("id")

	chatbotName := c.PostForm("chatbot_name")
	defaultMessage := c.PostForm("default_message")

	// Confirm ownership
	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2)`,
		chatbotID, userID).Scan(&exists)
	if err != nil || !exists {
		exceptions.Unauthorized("Invalid chatbot ID or access denied")
	}

	// Upload logo if present
	var logoPath string
	file, err := c.FormFile("logo")
	if err == nil {
		uploadDir := fmt.Sprintf("uploads/user_%d/chatbot_%s/logo", userID, chatbotID)
		_ = os.MkdirAll(uploadDir, os.ModePerm)

		filename := fmt.Sprintf("chatbot_logo_%s%s", chatbotID, filepath.Ext(file.Filename))
		logoPath = filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, logoPath); err != nil {
			exceptions.Internal("Logo upload failed")
		}
	} else {
		// Keep existing logo if no new one
		_ = db.DB.QueryRow(
			`SELECT logo_path FROM chatbots WHERE id = $1`,
			chatbotID,
		).Scan(&logoPath)
	}

	_, err = db.DB.Exec(`
		UPDATE chatbots
		SET chatbot_name = $1,
			logo_path = $2,
			default_message = $3
		WHERE id = $4 AND user_id = $5
	`,
		chatbotName, logoPath, defaultMessage, chatbotID, userID)

	if err != nil {
		exceptions.Internal("Failed to save settings")
	}

	responses.Success(c, "Settings updated successfully", gin.H{
		"logo_path": logoPath,
	})
}

func GetChatbotSettings(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotID := c.Param("id")

	// Confirm ownership
	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2)`,
		chatbotID, userID).Scan(&exists)
	if err != nil || !exists {
		exceptions.Unauthorized("Invalid chatbot ID or access denied")
	}

	var name, logo, message string
	err = db.DB.QueryRow(
		`SELECT chatbot_name, logo_path, default_message FROM chatbots WHERE id = $1`,
		chatbotID).Scan(&name, &logo, &message)

	if err != nil {
		// Return empty if not found
		responses.Success(c, "Chatbot settings fetched", gin.H{
			"chatbot_name":    "",
			"logo_path":       "",
			"default_message": "",
		})
		return
	}

	responses.Success(c, "Chatbot settings fetched", gin.H{
		"chatbot_name":    name,
		"logo_path":       logo,
		"default_message": message,
	})
}
