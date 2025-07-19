package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/globals"
	"github.com/stealcash/AgentFlow/app/logger"
	"github.com/stealcash/AgentFlow/app/middleware"
	"github.com/stealcash/AgentFlow/db"
	"github.com/stealcash/AgentFlow/routes"
	"log"
)

func main() {
	globals.RootDirPath = "./"
	globals.InitConfiguration()
	globals.LoadGeneralQuestionsConfig()

	// Main Db Connection
	if err := db.MainConnection(); err != nil {
		log.Fatal("Failed to connect to Main Db:", err)
	}
	defer db.CloseMainDb()

	// ElasticSearch connection
	if err := db.ElasticConnection(); err != nil {
		log.Fatal("Failed to connect to Elasticsearch:", err)
	}

	logger.InitLogger()
	defer logger.Close()

	// Gin router
	router := gin.Default()
	router.Use(middleware.RecoveryWithLogger())
	router.Use(middleware.CORS())

	// Proxies
	middleware.SetupTrustedProxies(router)
	router.MaxMultipartMemory = 8 << 20

	// Serve static
	router.Static("/uploads", "./uploads")
	rg := router.Group("/api")
	rg.Use(middleware.AuthMiddleware())

	routes.AuthRoutes(rg)
	routes.PublicRoutes(rg)
	routes.Profile(rg)
	routes.Chatbot(rg)
	routes.Domain(rg)
	routes.Plans(rg)

	log.Println("Running on " + globals.Config.App.Port)
	if err := router.Run(":" + globals.Config.App.Port); err != nil {
		log.Fatal(err)
	}
}
