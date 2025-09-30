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
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	ollamaService := services.NewOllamaService(cfg)
	similarityService := services.NewSimilarityService(ollamaService)
	handler := handlers.NewHandler(similarityService)

	router := routes.SetupRouter(handler)

	fmt.Println("Starting LLM Inference Service v1.2.0")
	fmt.Printf("Ollama URL: %s\n", cfg.OllamaURL)
	fmt.Printf("Model: %s\n", cfg.ModelName)
	fmt.Printf("Server starting on :%s\n", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
