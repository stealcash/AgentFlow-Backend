package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/logger"
	"net/http"
)

func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch ex := rec.(type) {
				case exceptions.HttpException:
					c.JSON(ex.StatusCode, gin.H{
						"status":  "error",
						"message": ex.Message,
					})
					c.Abort()

				case exceptions.HttpExceptionWithLog:
					logger.Error("Critical panic: %s | %s", ex.LogMessage, ex.Message)
					c.JSON(ex.StatusCode, gin.H{
						"status":  "error",
						"message": ex.Message,
					})
					c.Abort()

				case error:
					logger.Error("Unhandled panic: %v", ex)
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  "error",
						"message": "Internal Server Error",
					})
					c.Abort()

				default:
					logger.Error("Unknown panic: %v", rec)
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  "error",
						"message": "Internal Server Error",
					})
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
