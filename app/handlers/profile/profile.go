package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/db"
)

type UpdateProfileInput struct {
	CompanyName string `json:"company_name"`
}

func GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var email, company, userType string
	var parentID *int

	err := db.DB.QueryRow(`
		SELECT email, company_name, user_type, parent_id
		FROM users
		WHERE id = $1
	`, userID).Scan(&email, &company, &userType, &parentID)
	if err != nil {
		exceptions.Internal("Failed to fetch user")
	}

	rows, err := db.DB.Query(`
		SELECT id, chatbot_name, created_at
		FROM chatbots
		WHERE user_id = $1
	`, userID)
	if err != nil {
		exceptions.Internal("Failed to fetch chatbots")
	}
	defer rows.Close()

	var bots []gin.H
	for rows.Next() {
		var id int
		var name, createdAt string
		rows.Scan(&id, &name, &createdAt)
		bots = append(bots, gin.H{
			"id":         id,
			"name":       name,
			"created_at": createdAt,
		})
	}

	responses.Success(c, "Profile fetched successfully", gin.H{
		"user": gin.H{
			"email":     email,
			"company":   company,
			"user_type": userType,
			"parent_id": parentID,
			"chatbots":  bots,
		},
	})
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptions.BadRequest(err.Error())
	}

	_, err := db.DB.Exec(`UPDATE users SET company_name = $1 WHERE id = $2`, input.CompanyName, userID)
	if err != nil {
		exceptions.Internal("Failed to update profile")
	}

	responses.Success(c, "Profile updated successfully", nil)
}

func GetMe(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		exceptions.Unauthorized("Unauthorized")
	}

	var userType string
	err := db.DB.QueryRow(`SELECT user_type FROM users WHERE id = $1`, userID).Scan(&userType)
	if err != nil {
		exceptions.Internal("Could not find user")
	}

	responses.Success(c, "User fetched successfully", gin.H{
		"user_type": userType,
	})
}
