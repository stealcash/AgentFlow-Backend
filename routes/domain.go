package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/handlers/domain"
)

func Domain(rg *gin.RouterGroup) {
	{
		rgDomain := rg.Group("/v1/chatbots")
		rgDomain.GET("/:id/allowed-domains", domain.GetAllowedDomains)
		rgDomain.POST("/:id/allowed-domains", domain.AddAllowedDomain)

	}

}
