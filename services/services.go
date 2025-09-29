package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
func (s *OllamaService) CallAPI(prompt string) (string, error) {
	reqBody := models.OllamaRequest{
		Model:  s.config.ModelName,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(s.config.APIUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp models.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

// NERService handles Named Entity Recognition
type NERService struct {
	ollama *OllamaService
}

// NewNERService creates a new NERService instance
func NewNERService(ollama *OllamaService) *NERService {
	return &NERService{
		ollama: ollama,
	}
}

// ExtractEntities extracts named entities from text
func (s *NERService) ExtractEntities(text string) (*models.Entities, error) {
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
`, text)

	response, err := s.ollama.CallAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Parse the response
	data := safeJSONParse(response)

	// Post-processing: Convert to string slices and filter
	entities := &models.Entities{
		Persons:       filterByWordCount(toStringSlice(data["persons"]), 2),
		Locations:     toStringSlice(data["locations"]),
		Organizations: toStringSlice(data["organizations"]),
		Events:        filterByWordCount(toStringSlice(data["events"]), 2),
	}

	return entities, nil
}

// SentimentService handles sentiment analysis
type SentimentService struct {
	ollama *OllamaService
}

// NewSentimentService creates a new SentimentService instance
func NewSentimentService(ollama *OllamaService) *SentimentService {
	return &SentimentService{
		ollama: ollama,
	}
}

// AnalyzeSentiment analyzes the sentiment of text
func (s *SentimentService) AnalyzeSentiment(text string) (*models.Sentiment, error) {
	prompt := fmt.Sprintf(`
Analyze the sentiment of the following text. The text may be in Arabic, English, or German.
Return only valid JSON with this structure:
{
  "sentiment": "positive|negative|neutral",
  "confidence": 0.0-1.0
}

Text: %s
`, text)

	response, err := s.ollama.CallAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var sentimentData map[string]interface{}
	if err := json.Unmarshal([]byte(response), &sentimentData); err != nil {
		return &models.Sentiment{Sentiment: "neutral", Confidence: 0.0}, nil
	}

	// Extract sentiment and confidence
	sentiment := "neutral"
	confidence := 0.0

	if s, ok := sentimentData["sentiment"].(string); ok {
		sentiment = s
	}

	if conf, ok := sentimentData["confidence"].(float64); ok {
		confidence = conf
	} else if confStr, ok := sentimentData["confidence"].(string); ok {
		if parsed, err := strconv.ParseFloat(confStr, 64); err == nil {
			confidence = parsed
		}
	}

	return &models.Sentiment{
		Sentiment:  sentiment,
		Confidence: confidence,
	}, nil
}

// Helper functions

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
