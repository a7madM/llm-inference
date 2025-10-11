package services

import (
	"encoding/json"
	"fmt"
	"llm-inference/models"
	"strings"
)

// EntityEnhancementService handles entity quality improvement and deduplication
type EntityEnhancementService struct {
	ollama *OllamaService
}

// NewEntityEnhancementService creates a new EntityEnhancementService instance
func NewEntityEnhancementService(ollama *OllamaService) *EntityEnhancementService {
	return &EntityEnhancementService{ollama: ollama}
}

// EnhanceEntities processes an array of entities to improve quality and remove duplicates
func (s *EntityEnhancementService) EnhanceEntities(entities []string, entityType string) (*models.EntityEnhancementResponse, error) {
	if len(entities) == 0 {
		return &models.EntityEnhancementResponse{
			OriginalEntities: entities,
			EnhancedEntities: []string{},
			EntityType:       entityType,
			ProcessedCount:   0,
			RemovedCount:     0,
		}, nil
	}

	prompt := s.generateEnhancementPrompt(entities, entityType)

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
		OriginalEntities: entities,
		EnhancedEntities: enhancedResult.EnhancedEntities,
		EntityType:       entityType,
		ProcessedCount:   len(entities),
		RemovedCount:     len(entities) - len(enhancedResult.EnhancedEntities),
		Thinking:         response.Thinking,
	}, nil
}

// generateEnhancementPrompt creates the prompt for entity enhancement
func (s *EntityEnhancementService) generateEnhancementPrompt(entities []string, entityType string) string {
	entitiesStr := strings.Join(entities, "\", \"")

	prompt := fmt.Sprintf(`You are an expert in named entity recognition and data quality improvement. 

Your task is to enhance the quality of extracted %s entities by:
1. Correcting any obvious errors or typos
2. Standardizing formats and casing
5. Filtering out invalid or low-quality entries
6. Keeping only meaningful, well-formed entities
7. Keep entities in their original language without any translation
8. Return entities exactly as they appear in their source language


Input entities: ["%s"]

Rules:
1. Return only valid, high-quality %s entities
2. Remove duplicates and near-duplicates (e.g., "John Smith" and "john smith" are the same)
3. Fix obvious typos and formatting issues
4. Remove entries that are clearly not valid %s entities
5. Standardize capitalization appropriately for %s entities
6. Removing duplicates, it means return only one instance of each entity, e.g., "John Smith" and "john smith" are the same
7. Consider arabic letters, ي or ى may be the same letter
8. Never translate entities - keep them in their original language
9. Return results in JSON format only

Required JSON structure:
{
  "enhanced_entities": ["entity1", "entity2", "entity3"]
}

Do not include any explanations, markdown, or extra text - only the JSON response.`,
		entityType, entitiesStr, entityType, entityType, entityType)

	return prompt
}

// parseEnhancementResponse parses the LLM's JSON response
func (s *EntityEnhancementService) parseEnhancementResponse(jsonResponse string) (*EntityEnhancementResult, error) {
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
func (s *EntityEnhancementService) extractCleanJSON(content string) string {
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")

	if start != -1 && end != -1 && start < end {
		return content[start : end+1]
	}

	return ""
}

// EntityEnhancementResult represents the parsed result from LLM
type EntityEnhancementResult struct {
	EnhancedEntities []string `json:"enhanced_entities"`
}

// Legacy methods for backward compatibility

func (s *OllamaService) NameEntityEnhance(text1, text2, entityType string) (models.OllamaResponse, error) {
	return s.CallAPI(generatePrompt(text1, text2, entityType))
}
