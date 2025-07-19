package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/handlers/plans"
)

func Plans(rg *gin.RouterGroup) {
	{
		rgPlan := rg.Group("/v1/plans")
		rgPlan.POST("", plans.CreatePlan)
		rgPlan.GET("", plans.ListPlans)
	}

	{
		rgSubscriber := rg.Group("/v1/subscription")
		rgSubscriber.POST("", plans.SubscribeUser)
		rgSubscriber.GET("", plans.GetSubscription)
	}

}
