package services

import (
	"encoding/json"
	"fmt"
	"llm-inference/models"
	"strings"
)

// EntityVerifier handles entity quality improvement and deduplication
type EntityVerifier struct {
	ollama *OllamaService
}

// NewEntityEnhancementService creates a new EntityVerifier instance
func NewEntityEnhancementService(ollama *OllamaService) *EntityVerifier {
	return &EntityVerifier{ollama: ollama}
}

// EnhanceEntities processes an array of entities to improve quality and remove duplicates
func (s *EntityVerifier) Verify(entity, entityType string) (*models.EntityEnhancementResponse, error) {
	if entity == "" {
		return &models.EntityEnhancementResponse{
			Entity:   entity,
			Verified: false,
		}, nil
	}

	prompt := s.generateEnhancementPrompt(entity, entityType)

	response, err := s.ollama.CallAPI(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance entities: %v", err)
	}

	// Parse the LLM response
	response.ParseJSON()
	enhancedResult, err := s.parseEnhancementResponse(response.JSONResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse enhancement response: %v", err)
	}

	return &models.EntityEnhancementResponse{
		Entity:   entity,
		Verified: enhancedResult.Verified,
	}, nil
}

// generateEnhancementPrompt creates the prompt for entity enhancement
func (s *EntityVerifier) generateEnhancementPrompt(entity, entityType string) string {
	prompt := fmt.Sprintf(`You are an expert in named entity recognition and data quality improvement.

Your task is to enhance the quality of extracted %s entities by:
1. Determining if each entity is valid and meaningful
3. Never translating entities - keep them in their original language
4. Considering Arabic letters, ي or ى may be the same letter
5. Returning results in JSON format only

Input entity: "%s"

Rules:
1. For the entity, determine if it is valid and meaningful
2. Never translate entities - keep them in their original language
3. Consider arabic letters, ي or ى may be the same letter
4. Consider all forms of the entity, including singular/plural and different grammatical cases
5. Returns verified as true if the entity is valid, otherwise false
6. The entity is not valid if it is too generic, vague, or does not provide useful information
7. The entity is not a valid location if it is a non-specific place like "city", "country", or "region"
8. The entity is not a valid person name if it is a common noun or title like "doctor", "engineer", or "teacher"
9. The entity is not a valid organization if it is a generic term like "company", "institution", or "agency"
10. The entity is not valid location if it contains two or more different locations combined, like "الولايات المتحدة وكندا"
11. Return results in JSON format only

Required JSON structure:
{
  "entity": entity_name,
  "verified": true/false
}
`, entityType, entity)

	return prompt
}

// parseEnhancementResponse parses the LLM's JSON response
func (s *EntityVerifier) parseEnhancementResponse(jsonResponse string) (*EntityEnhancementResult, error) {
	var result EntityEnhancementResult

	if err := json.Unmarshal([]byte(jsonResponse), &result); err != nil {
		// Try to extract JSON if it's embedded in other text
		if cleanJSON := s.extractCleanJSON(jsonResponse); cleanJSON != "" {
			if err := json.Unmarshal([]byte(cleanJSON), &result); err == nil {
				return &result, nil
			}
		}
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return &result, nil
}

// extractCleanJSON attempts to extract valid JSON from mixed content
func (s *EntityVerifier) extractCleanJSON(content string) string {
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")

	if start != -1 && end != -1 && start < end {
		return content[start : end+1]
	}

	return ""
}

// EntityEnhancementResult represents the parsed result from LLM
type EntityEnhancementResult struct {
	Entity   string `json:"entity"`
	Verified bool   `json:"verified"`
}

// Legacy methods for backward compatibility

func (s *OllamaService) NameEntityEnhance(text1, text2, entityType string) (models.OllamaResponse, error) {
	return s.CallAPI(generatePrompt(text1, text2, entityType))
}
