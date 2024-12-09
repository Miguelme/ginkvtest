package main

import (
	"ginkvtest/internal/router"
	"ginkvtest/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// Initialize Aurora service

	auroraService, err := service.NewAuroraService()
	if err != nil {
		log.Fatalf("Failed to initialize Aurora service: %v", err)
	}

	// Initialize DynamoDB service
	dynamoService := service.NewDynamoService()

	// Initialize Redis service
	redisService, err := service.NewRedisService()
	if err != nil {
		log.Fatalf("Failed to initialize Redis service: %v", err)
	}

	// Set up the router
	r := gin.Default()
	router.SetupRoutes(r, auroraService, dynamoService, redisService)

	// Start the server
	log.Println("Server running on :8080")
	r.Run(":8080")
}
