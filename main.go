package main

import (
	"log"
	"os"

	"blog-api/config"
	"blog-api/routes"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Connect to database
	config.ConnectDatabase()

	// Setup routes
	router := routes.SetupRoutes(logger)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.Info("Starting server", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
