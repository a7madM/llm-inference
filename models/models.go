package models

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
}

// OllamaResponse represents a response from the Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
}

type SimilarityResponse struct {
	Text1           string  `json:"text1"`
	Text2           string  `json:"text2"`
	SimilarityScore float64 `json:"similarity_score"`
	ShouldMerge     bool    `json:"should_be_merged"`
}
