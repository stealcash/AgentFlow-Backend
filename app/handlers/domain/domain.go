package domain

import (
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/db"

	"github.com/gin-gonic/gin"
)

func AddAllowedDomain(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotID := c.Param("id")

	var input struct {
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptions.BadRequest(err.Error())
	}

	var exists bool
	err := db.DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2)`,
		chatbotID, userID).Scan(&exists)
	if err != nil {
		exceptions.Internal("Database error")
	}
	if !exists {
		exceptions.Forbidden("Unauthorized")
	}

	_, err = db.DB.Exec(
		`INSERT INTO allowed_domains (chatbot_id, domain) VALUES ($1, $2)`,
		chatbotID, input.Domain)
	if err != nil {
		exceptions.Internal("Failed to add domain")
	}

	responses.Success(c, "Domain added", nil)
}

func GetAllowedDomains(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotID := c.Param("id")

	var exists bool
	err := db.DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2)`,
		chatbotID, userID).Scan(&exists)
	if err != nil {
		exceptions.Internal("Database error")
	}
	if !exists {
		exceptions.Forbidden("Unauthorized")
	}

	rows, err := db.DB.Query(
		`SELECT id, domain FROM allowed_domains WHERE chatbot_id = $1`,
		chatbotID)
	if err != nil {
		exceptions.Internal("Failed to fetch domains")
	}
	defer rows.Close()

	var domains []gin.H
	for rows.Next() {
		var id int
		var domain string
		rows.Scan(&id, &domain)
		domains = append(domains, gin.H{
			"id":     id,
			"domain": domain,
		})
	}

	responses.Success(c, "Domains fetched", gin.H{"domains": domains})
}
