package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/handlers/auth"
)

func AuthRoutes(rg *gin.RouterGroup) {
	{
		rgAuth := rg.Group("/v1/auth")
		rgAuth.POST("/signup", auth.Signup)
		rgAuth.POST("/login", auth.Login)
	}
}
