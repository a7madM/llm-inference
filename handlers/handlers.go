package handlers

import (
	"llm-inference/services"
)

// Handler holds the service dependencies
type Handler struct {
	similarityService *services.SimilarityService
	entityVerifier    *services.EntityVerifier
}

func NewHandler(similarityService *services.SimilarityService, entityVerifier *services.EntityVerifier) *Handler {
	return &Handler{
		similarityService: similarityService,
		entityVerifier:    entityVerifier,
	}
}
