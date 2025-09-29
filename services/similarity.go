package services

import (
	"encoding/json"
	"fmt"
	"llm-inference/models"
	"strings"
)

func (s *OllamaService) Similarity(text1, text2, entityType string) (models.OllamaResponse, error) {
	prompt := fmt.Sprintf(`You are an intilligent model. Given the following two texts
	 compute a similarity score between 0 and 1, where 1 means identical meaning and 0 means completely different. Also, indicate if the texts should be merged based on their similarity. 
	 

	Rules:
	1. These two texts are from entity type: %s.
	2. Don't compare the texts based on exact wording; focus on meaning.
	3. Analyze the semantic meaning of both texts.
	4. Consider context, synonyms, and phrasing.
	5. Return only valid JSON with this structure:
	{
		"similarity_score": float,
		"should_be_merged": boolean
	}
	6. Do not translate the texts.
	7. Keep the response strictly to the JSON format.
	8. Do not include any explanations, extra text, or markdown.
	9. Round the similarity score to two decimal places.

	Text 1: %s\nText 2: %s`, entityType, text1, text2)
	return s.CallAPI(prompt)
}

type SimilarityService struct {
	ollama *OllamaService
}

func NewSimilarityService(ollama *OllamaService) *SimilarityService {
	return &SimilarityService{
		ollama: ollama,
	}
}

func (s *SimilarityService) ComputeSimilarity(text1, text2, entityType string) (*models.SimilarityResponse, error) {
	response, err := s.ollama.Similarity(text1, text2, entityType)

	if err != nil {
		return nil, err
	}

	data := s.ParseJSON(response.Response)

	return &models.SimilarityResponse{
		Text1:           text1,
		Text2:           text2,
		SimilarityScore: data.SimilarityScore,
		ShouldMerge:     data.ShouldBeMerged,
	}, nil
}

// Safe JSON parsing with fallback
func (s *SimilarityService) ParseJSON(rawOutput string) ResponseSchema {

	var result ResponseSchema

	start := strings.Index(rawOutput, "{")
	end := strings.LastIndex(rawOutput, "}")
	if start != -1 && end != -1 && start < end {
		jsonPart := rawOutput[start : end+1]
		err := json.Unmarshal([]byte(jsonPart), &result)
		if err == nil {
			return result
		}
	}

	return result
}

type ResponseSchema struct {
	SimilarityScore float64 `json:"similarity_score"`
	ShouldBeMerged  bool    `json:"should_be_merged"`
}
