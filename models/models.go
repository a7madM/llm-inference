package models

import (
	"strings"
)

// InputText represents the input for API requests
type InputText struct {
	Text string `json:"text" binding:"required"`
}

// Entities represents the named entities extraction response
type Entities struct {
	Persons       []string `json:"persons"`
	Locations     []string `json:"locations"`
	Organizations []string `json:"organizations"`
	Events        []string `json:"events"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

// ServiceInfo represents the service information response
type ServiceInfo struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Endpoints   map[string]string `json:"endpoints"`
}

// OllamaRequest represents a request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Token  string `json:"token,omitempty"`
}

// OllamaResponse represents a response from the Ollama API
type OllamaResponse struct {
	Response     string `json:"response"`
	Thinking     string `json:"thinking"`
	JSONResponse string `json:"json_response"`
}

func (o *OllamaResponse) ParseJSON() {
	rawOutput := o.Response

	// Find the first '{' and the last '}' to extract the JSON part

	start := strings.Index(rawOutput, "{")
	end := strings.LastIndex(rawOutput, "}")

	if start != -1 {
		o.Thinking = rawOutput[0:start] // Text before JSON
	} else {
		o.Thinking = rawOutput // No JSON found, treat all as thinking
	}

	if start != -1 && end != -1 && start < end {
		jsonPart := rawOutput[start : end+1]
		o.JSONResponse = jsonPart
	}
}

type SimilarityResponse struct {
	Text1           string  `json:"text1"`
	Text2           string  `json:"text2"`
	SimilarityScore float64 `json:"similarity_score"`
	ShouldMerge     bool    `json:"should_be_merged"`
	Thinking        string  `json:"thinking,omitempty"`
}

// EntityVerifierRequest represents the request for entity verification
type EntityVerifierRequest struct {
	Entity string `json:"entity"`
	Type   string `json:"type"`
}

// EntityVerifierResponse represents the verified entities response
type EntityVerifierResponse struct {
	Entity   string `json:"entity"`
	Verified bool   `json:"verified"`
}
