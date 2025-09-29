package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	OllamaURL string
	ModelName string
	APIUrl    string
	Port      string
	GinMode   string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load environment variables from .env file if it exists
	godotenv.Load()

	cfg := &Config{
		OllamaURL: getEnv("OLLAMA_URL", "http://localhost:11434"),
		ModelName: getEnv("MODEL_NAME", "deepseek-r1:1.5b"),
		Port:      getEnv("PORT", "8090"),
		GinMode:   getEnv("GIN_MODE", "release"),
	}
	cfg.APIUrl = cfg.OllamaURL + "/api/generate"

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
