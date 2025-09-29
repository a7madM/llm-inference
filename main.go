package main

import (
	"fmt"
	"log"

	"llm-inference/config"
	"llm-inference/handlers"
	"llm-inference/routes"
	"llm-inference/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize services
	ollamaService := services.NewOllamaService(cfg)
	nerService := services.NewNERService(ollamaService)
	sentimentService := services.NewSentimentService(ollamaService)

	// Initialize handlers
	handler := handlers.NewHandler(nerService, sentimentService)

	// Setup routes
	router := routes.SetupRouter(handler)

	// Print startup information
	fmt.Printf("Starting Local Multilingual NER Service v1.2.0\n")
	fmt.Printf("Ollama URL: %s\n", cfg.OllamaURL)
	fmt.Printf("Model: %s\n", cfg.ModelName)
	fmt.Printf("Server starting on :%s\n", cfg.Port)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
