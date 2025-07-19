package middleware

import (
	"github.com/stealcash/AgentFlow/app"
	"github.com/stealcash/AgentFlow/app/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlFullPath := c.FullPath()
		if strings.HasPrefix(urlFullPath, "/uploads/") ||
			strings.HasPrefix(urlFullPath, "/api/v1/auth/") ||
			strings.HasPrefix(urlFullPath, "/api/v1/public") {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		// Strip "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userID := int(claims["user_id"].(float64))
		userType := claims["user_type"].(string)

		// Save into context
		c.Set("user_id", userID)
		c.Set("user_type", userType)
		app.SetCurrentUserInReqContext(c, claims)

		c.Next()
	}
}
