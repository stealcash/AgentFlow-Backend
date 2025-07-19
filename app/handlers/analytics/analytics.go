package analytics

import (
	"database/sql"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/db"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsInput struct {
	ChatbotID      int    `json:"chatbot_id" binding:"required"`
	Domain         string `json:"domain"`
	CategoryID     *int   `json:"category_id"`
	QuestionID     *int   `json:"question_id"`
	InputQuery     string `json:"input_query"`
	ResponseSource string `json:"response_source"`
}

func LogAnalytics(chatbotID int, domain string, categoryID, questionID *int, inputQuery string, responseSource string) error {
	_, err := db.DB.Exec(`
		INSERT INTO analytics (
			chatbot_id, domain, category_id, question_id,
			input_query, response_source, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		chatbotID, domain, categoryID, questionID,
		inputQuery, responseSource, time.Now(),
	)
	return err
}
func GetAnalytics(c *gin.Context) {
	chatbotID := c.Param("id")

	var totalVisits, totalChats, messagesSent int
	var createdAt sql.NullTime

	err := db.DB.QueryRow(`
		SELECT 
			COUNT(DISTINCT domain) as total_visits,
			COUNT(DISTINCT input_query) as total_chats,
			COUNT(*) as messages_sent,
			MIN(created_at) as created_at
		FROM analytics
		WHERE chatbot_id = $1
	`, chatbotID).Scan(&totalVisits, &totalChats, &messagesSent, &createdAt)

	if err != nil {
		exceptions.Internal("Failed to aggregate analytics")
	}

	// fallback if no records yet
	createdAtValue := time.Now()
	if createdAt.Valid {
		createdAtValue = createdAt.Time
	}

	// Fetch logs...
	rows, err := db.DB.Query(`
		SELECT id, domain, input_query, response_source, category_id, question_id, created_at
		FROM analytics
		WHERE chatbot_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, chatbotID)

	if err != nil {
		exceptions.Internal("Database error while fetching analytics logs")
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int
		var domain, inputQuery, source string
		var categoryId, questionId sql.NullInt64
		var createdAt time.Time

		rows.Scan(&id, &domain, &inputQuery, &source, &categoryId, &questionId, &createdAt)

		logs = append(logs, gin.H{
			"id":              id,
			"domain":          domain,
			"input_query":     inputQuery,
			"response_source": source,
			"category_id":     ifNullInt(categoryId),
			"question_id":     ifNullInt(questionId),
			"created_at":      createdAt,
		})
	}

	responses.Success(c, "Analytics fetched", gin.H{
		"analytics": gin.H{
			"total_visits":  totalVisits,
			"total_chats":   totalChats,
			"messages_sent": messagesSent,
			"created_at":    createdAtValue,
		},
		"logs": logs,
	})
}

func ifNullInt(val sql.NullInt64) interface{} {
	if val.Valid {
		return val.Int64
	}
	return nil
}
