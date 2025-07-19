package script

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/app/globals"
	"github.com/stealcash/AgentFlow/db"
	"strconv"
)

// generate a secure random API key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateEmbed(c *gin.Context) {
	userID := c.GetInt("user_id")

	chatbotIDStr := c.Param("id")
	chatbotID, err := strconv.Atoi(chatbotIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid chatbot ID")
		return
	}

	// Check for optional `regen` flag
	regen := c.DefaultQuery("regenerate", "false") == "true"

	var publicAPIKey string
	err = db.DB.QueryRow(
		`SELECT public_api_key FROM chatbots WHERE id = $1 AND user_id = $2`,
		chatbotID, userID).Scan(&publicAPIKey)
	if err != nil {
		exceptions.Forbidden("Chatbot not found or unauthorized")
		return
	}

	if regen {
		newKey, err := GenerateAPIKey()
		if err != nil {
			exceptions.Internal("Failed to generate new API key")
			return
		}

		_, err = db.DB.Exec(
			`UPDATE chatbots SET public_api_key = $1 WHERE id = $2`,
			newKey, chatbotID)
		if err != nil {
			exceptions.Internal("Failed to update API key")
			return
		}

		publicAPIKey = newKey
	}

	config := map[string]string{
		"apiKey": publicAPIKey,
		"apiUrl": globals.Config.App.ApiUrl,
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		exceptions.Internal("Failed to build config JSON")
		return
	}

	configBase64 := base64.StdEncoding.EncodeToString(configJSON)

	embedScript := fmt.Sprintf(
		`<script src="%s/scripts/chatbot.js" data-config="%s"></script>`,
		globals.Config.FrontEnd.Path, configBase64,
	)

	responses.Success(c, "Embed script generated", gin.H{
		"embed_script":   embedScript,
		"public_api_key": publicAPIKey,
		"regenerated":    regen,
	})
}
