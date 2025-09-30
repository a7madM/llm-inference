package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"llm-inference/config"
	"llm-inference/models"
)

// OllamaService handles communication with Ollama API
type OllamaService struct {
	config *config.Config
}

// NewOllamaService creates a new OllamaService instance
func NewOllamaService(cfg *config.Config) *OllamaService {
	return &OllamaService{
		config: cfg,
	}
}

// CallAPI makes a request to the Ollama API
func (s *OllamaService) CallAPI(prompt string) (models.OllamaResponse, error) {
	reqBody := models.OllamaRequest{
		Model:  s.config.ModelName,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(s.config.APIUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.OllamaResponse{}, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp models.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}
	ollamaResp.ParseJSON()

	fmt.Println("Ollama API Response:", ollamaResp)
	return ollamaResp, nil
}
