package middleware

import (
	"github.com/stealcash/AgentFlow/app/globals"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupTrustedProxies(router *gin.Engine) {
	err := router.SetTrustedProxies(globals.Config.TrustedProxies)
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}
}
