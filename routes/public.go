package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/handlers/public"
)

func PublicRoutes(rg *gin.RouterGroup) {
	{
		rgPublicRoutes := rg.Group("/v1/public")
		rgPublicRoutes.GET("/chatbot-profile", public.GetPublicChatbotProfile)
		rgPublicRoutes.POST("/chatbot", public.PublicChatQuery)
	}
}
