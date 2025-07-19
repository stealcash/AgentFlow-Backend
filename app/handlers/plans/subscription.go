package plans

import (
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PlanInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Features    string  `json:"features"` // JSON string or comma-separated
	Price       float64 `json:"price"`
}

// CreatePlan inserts a new plan only for superadmin
func CreatePlan(c *gin.Context) {
	var input PlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO plans (name, description, features, price)
		VALUES ($1, $2, $3, $4)
	`, input.Name, input.Description, input.Features, input.Price)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}
	responses.Success(c, "Plan created", gin.H{"message": "Plan created"})

}

// ListPlans returns all available plans
func ListPlans(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, name, description, features, price
		FROM plans
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch plans"})
		return
	}
	defer rows.Close()

	var plans []gin.H
	for rows.Next() {
		var id int
		var name, desc, features string
		var price float64

		if err := rows.Scan(&id, &name, &desc, &features, &price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan failed"})
			return
		}

		plans = append(plans, gin.H{
			"id":          id,
			"name":        name,
			"description": desc,
			"features":    features,
			"price":       price,
		})
	}

	if plans == nil {
		plans = []gin.H{}
	}

	responses.Success(c, "List of plan successful", gin.H{"plans": plans})
}

type SubscriptionInput struct {
	PlanID int `json:"plan_id" binding:"required"`
}

// SubscribeUser starts a subscription for the authenticated user
func SubscribeUser(c *gin.Context) {
	userID := c.GetInt("user_id")

	var input SubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO subscriptions (user_id, plan_id, start_date, end_date, status)
		VALUES ($1, $2, CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', 'active')
	`, userID, input.PlanID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Subscription failed"})
		return
	}
	responses.Success(c, "Subscription successful", gin.H{"message": "Subscription successful"})

}

// GetSubscription returns the latest subscription for the user
func GetSubscription(c *gin.Context) {
	userID := c.GetInt("user_id")

	row := db.DB.QueryRow(`
		SELECT s.id, p.name, p.description, p.features, p.price,
		       s.start_date, s.end_date, s.status
		FROM subscriptions s
		JOIN plans p ON s.plan_id = p.id
		WHERE s.user_id = $1
		ORDER BY s.created_at DESC
		LIMIT 1
	`, userID)

	var id int
	var planName, desc, features, status string
	var price float64
	var startDate, endDate time.Time

	err := row.Scan(&id, &planName, &desc, &features, &price, &startDate, &endDate, &status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"subscription": nil})
		return
	}
	responses.Success(c, "Subscription successful", gin.H{
		"subscription": gin.H{
			"plan": gin.H{
				"name":        planName,
				"description": desc,
				"features":    features,
				"price":       price,
			},
			"start_date": startDate,
			"end_date":   endDate,
			"status":     status,
		},
	})

}
