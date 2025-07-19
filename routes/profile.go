package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/handlers/profile"
)

func Profile(rg *gin.RouterGroup) {
	{
		rgProfile := rg.Group("/v1/profile")
		rgProfile.GET("", profile.GetProfile)
		rgProfile.PUT("", profile.UpdateProfile)
		rgProfile.GET("/me", profile.GetMe)

	}

}
