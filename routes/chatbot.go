package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/handlers/analytics"
	"github.com/stealcash/AgentFlow/app/handlers/chatbot"
	"github.com/stealcash/AgentFlow/app/handlers/script"
)

func Chatbot(rg *gin.RouterGroup) {
	{
		rgChatbot := rg.Group("/v1/chatbots")
		rgChatbot.POST("", chatbot.CreateChatbot)
		rgChatbot.GET("", chatbot.ListChatbots)
		rgChatbot.DELETE("/:id", chatbot.DeleteChatbot)

		rgChatbot.GET("/:id/analytics", analytics.GetAnalytics)
		rgChatbot.POST("/:id/embed", script.GenerateEmbed)

		rgChatbot.POST("/:id/general-questions", chatbot.AddGeneralQuestion)
		rgChatbot.GET("/:id/general-questions", chatbot.GetGeneralQuestions)
		rgChatbot.DELETE("/:id/general-questions/:question_id", chatbot.DeleteGeneralQuestion)

		rgChatbot.POST("/:id/questions", chatbot.AddQuestion)
		rgChatbot.GET("/:id/questions/:category_id", chatbot.GetQuestionsByCategory)
		rgChatbot.DELETE("/questions/:id", chatbot.DeleteQuestion)

	}
	{
		rgChatSetting := rg.Group("/v1/chatbots/:id/settings")
		rgChatSetting.POST("", chatbot.UploadChatbotSettings)
		rgChatSetting.GET("", chatbot.GetChatbotSettings)
	}

	{
		rgChatSetting := rg.Group("/v1/chatbots/:id/categories")
		rgChatSetting.POST("", chatbot.CreateCategory)
		rgChatSetting.GET("", chatbot.GetCategories)
		rgChatSetting.DELETE("/:cat_id", chatbot.DeleteCategory)
		rgChatSetting.POST("/image", chatbot.UploadCategoryImage) // not working

	}

}
