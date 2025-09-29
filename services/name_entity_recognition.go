package services

import (
	"encoding/json"
	"fmt"
	"llm-inference/models"
	"strings"
)

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
	data := s.ParseJSON(response.Response)

	// Post-processing: Convert to string slices and filter
	entities := &models.Entities{
		Persons:       filterByWordCount(toStringSlice(data["persons"]), 2),
		Locations:     toStringSlice(data["locations"]),
		Organizations: toStringSlice(data["organizations"]),
		Events:        filterByWordCount(toStringSlice(data["events"]), 2),
	}

	return entities, nil
}

// Helper functions

// Safe JSON parsing with fallback
func (s *NERService) ParseJSON(rawOutput string) map[string]interface{} {
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
