package main

import (
	"testing"

	"llm-inference/config"
	"llm-inference/handlers"
	"llm-inference/services"
)

func TestServiceInitialization(t *testing.T) {
	// Test config loading
	cfg := config.Load()
	if cfg == nil {
		t.Fatal("Config should not be nil")
	}

	if cfg.OllamaURL == "" {
		t.Fatal("OllamaURL should not be empty")
	}

	if cfg.ModelName == "" {
		t.Fatal("ModelName should not be empty")
	}

	// Test service initialization
	ollamaService := services.NewOllamaService(cfg)
	if ollamaService == nil {
		t.Fatal("OllamaService should not be nil")
	}

	nerService := services.NewNERService(ollamaService)
	if nerService == nil {
		t.Fatal("NERService should not be nil")
	}

	sentimentService := services.NewSentimentService(ollamaService)
	if sentimentService == nil {
		t.Fatal("SentimentService should not be nil")
	}

	// Test handler initialization
	handler := handlers.NewHandler(nerService, sentimentService)
	if handler == nil {
		t.Fatal("Handler should not be nil")
	}
}
