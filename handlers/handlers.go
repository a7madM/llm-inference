package handlers

import (
	"fmt"
	"time"

	"llm-inference/models"
	"llm-inference/services"

	"github.com/gofiber/fiber/v2"
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

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	response := models.HealthResponse{
		Status:  "healthy",
		Service: "LLM Inference Service",
		Uptime:  time.Now().UTC().Format(time.RFC3339),
		Version: "0.1.0",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) ComputeSimilarity(c *fiber.Ctx) error {
	text1 := c.Query("text1")
	text2 := c.Query("text2")
	entityType := c.Query("entity_type")
	if text1 == "" || text2 == "" || entityType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Both text1 and text2, and EntityType query parameters are required"})
	}

	result, err := h.similarityService.ComputeSimilarity(text1, text2, entityType)
	if err != nil {
		fmt.Println("Error computing similarity:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to compute similarity"})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Verifier handles the entity verification endpoint
func (h *Handler) Verifier(c *fiber.Ctx) error {
	var input models.EntityEnhancementRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// Log request
	fmt.Printf("Enhancing %s with type %s\n", input.Entity, input.Type)

	startTime := time.Now()

	// Verify entity
	result, err := h.entityVerifier.Verify(input.Entity, input.Type)
	if err != nil {
		fmt.Printf("Error enhancing entities: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Failed to enhance entities"})
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Entity enhancement completed in %.2f seconds.\n", elapsed.Seconds())

	return c.Status(fiber.StatusOK).JSON(result)
}

// ServiceInfo handles the service information endpoint
func (h *Handler) ServiceInfo(c *fiber.Ctx) error {
	response := models.ServiceInfo{
		Title:       "LLM Inference Service",
		Description: "A service for inference tasks using large language models.",
		Version:     "0.1.0",
		Endpoints: map[string]string{
			"similarity": "/api/v1/similarity",
			"verify":     "/api/v1/verify",
			"health":     "/health",
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
