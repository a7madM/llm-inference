package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Configuration
type Config struct {
	OllamaURL string
	ModelName string
	APIUrl    string
}

// Input/Output Models
type InputText struct {
	Text string `json:"text" binding:"required"`
}

type Entities struct {
	Persons       []string `json:"persons"`
	Locations     []string `json:"locations"`
	Organizations []string `json:"organizations"`
	Events        []string `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Ollama API request/response structures
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

// Global configuration
var config Config

func init() {
	// Load environment variables from .env file if it exists
	godotenv.Load()

	config = Config{
		OllamaURL: getEnv("OLLAMA_URL", getEnv("OLLAMA_URL", "localhost:11434")),
		ModelName: getEnv("MODEL_NAME", getEnv("MODEL_NAME", "deepseek-r1:1.5b")),
	}
	config.APIUrl = config.OllamaURL + "/api/generate"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Safe JSON parsing with fallback
func safeJSONParse(rawOutput string) map[string]interface{} {
	var result map[string]interface{}

	// Try to parse the entire string first
	if err := json.Unmarshal([]byte(rawOutput), &result); err == nil {
		return result
	}

	// Try to extract JSON from between first { and last }
	start := strings.Index(rawOutput, "{")
	end := strings.LastIndex(rawOutput, "}")

	if start != -1 && end != -1 && end > start {
		jsonStr := rawOutput[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
			return result
		}
	}

	// Return empty schema as fallback
	return map[string]interface{}{
		"persons":       []string{},
		"locations":     []string{},
		"organizations": []string{},
		"events":        []string{},
	}
}

// Helper function to convert interface{} to []string
func toStringSlice(data interface{}) []string {
	if data == nil {
		return []string{}
	}

	switch v := data.(type) {
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	case []string:
		return v
	default:
		return []string{}
	}
}

// Count words in a string
func countWords(text string) int {
	return len(strings.Fields(strings.TrimSpace(text)))
}

// Filter strings by minimum word count
func filterByWordCount(items []string, minWords int) []string {
	var filtered []string
	for _, item := range items {
		if countWords(item) >= minWords {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// Call Ollama API
func callOllamaAPI(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  config.ModelName,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(config.APIUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

// Extract entities endpoint
func extractEntities(c *gin.Context) {
	var input InputText
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Log request
	textPreview := input.Text
	if len(textPreview) > 50 {
		textPreview = textPreview[:50] + "..."
	}
	fmt.Printf("Received text for NER: %s\n", textPreview)

	startTime := time.Now()

	prompt := fmt.Sprintf(`
		You are an information extraction system. 
		Extract named entities from the following text **without translating it**.
		The text may be in Arabic, English, or German.

		Rules:
		1. Person names must consist of at least **two words**.
		2. Events must be **phrases**, not single words.
		3. Return only valid JSON with this structure:
		{
		"persons": ["string"],
		"locations": ["string"],
		"organizations": ["string"],
		"events": ["string"]
		}
		4. Do not translate the text. Do not include explanations, extra text, or markdown. 
		5. Keep all names and words in the original language.

		Text: %s
		`, input.Text)

	// Call Ollama API
	response, err := callOllamaAPI(prompt)
	if err != nil {
		log.Printf("Error calling Ollama API: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Unable to reach Ollama API"})
		return
	}

	// Parse the response
	data := safeJSONParse(response)

	// Post-processing: Convert to string slices and filter
	entities := Entities{
		Persons:       filterByWordCount(toStringSlice(data["persons"]), 2),
		Locations:     toStringSlice(data["locations"]),
		Organizations: toStringSlice(data["organizations"]),
		Events:        filterByWordCount(toStringSlice(data["events"]), 2),
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Extraction completed in %.2f seconds.\n", elapsed.Seconds())

	c.JSON(http.StatusOK, entities)
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "Local Multilingual NER Service",
		"uptime":  time.Now().UTC().Format(time.RFC3339),
		"version": "1.2.0",
	})
}

func setupRouter() *gin.Engine {
	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API routes
	v1 := router.Group("/api/v1")
	v1.POST("/entities", extractEntities)

	// Health check
	router.GET("/health", healthCheck)

	// Root endpoint with service info
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title":       "Local Multilingual NER Service",
			"description": "Extracts Persons, Locations, Organizations, Events from text (Arabic, English, German). Preserves original language.",
			"version":     "1.2.0",
			"endpoints": gin.H{
				"entities":  "/api/v1/entities",
				"sentiment": "/api/v1/sentiment",
				"health":    "/health",
			},
		})
	})

	return router
}

func main() {
	fmt.Printf("Starting Local Multilingual NER Service v1.2.0\n")
	fmt.Printf("Ollama URL: %s\n", config.OllamaURL)
	fmt.Printf("Model: %s\n", config.ModelName)

	router := setupRouter()

	// Start server
	fmt.Printf("Server starting on :%d\n", 8090)
	if err := router.Run(":8090"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
