package handlers

import (
	"fmt"
	"net/http"
	"time"

	"llm-inference/models"
	"llm-inference/services"

	"github.com/gin-gonic/gin"
)

// Handler holds the service dependencies
type Handler struct {
	similarityService        *services.SimilarityService
	entityEnhancementService *services.EntityEnhancementService
}

func NewHandler(similarityService *services.SimilarityService, entityEnhancementService *services.EntityEnhancementService) *Handler {
	return &Handler{
		similarityService:        similarityService,
		entityEnhancementService: entityEnhancementService,
	}
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	response := models.HealthResponse{
		Status:  "healthy",
		Service: "LLM Inference Service",
		Uptime:  time.Now().UTC().Format(time.RFC3339),
		Version: "0.1.0",
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) ComputeSimilarity(c *gin.Context) {
	text1 := c.Query("text1")
	text2 := c.Query("text2")
	entityType := c.Query("entity_type")
	if text1 == "" || text2 == "" || entityType == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Both text1 and text2, and EntityType query parameters are required"})
		return
	}

	result, err := h.similarityService.ComputeSimilarity(text1, text2, entityType)

	if err != nil {
		fmt.Println("Error computing similarity:", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to compute similarity"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// EnhanceEntities handles the entity enhancement endpoint
func (h *Handler) EnhanceEntities(c *gin.Context) {
	var input models.EntityEnhancementRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Log request
	fmt.Printf("Enhancing %d %s entities...\n", len(input.Entities), input.EntityType)

	startTime := time.Now()

	// Enhance entities
	result, err := h.entityEnhancementService.EnhanceEntities(input.Entities, input.EntityType)
	if err != nil {
		fmt.Printf("Error enhancing entities: %v\n", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to enhance entities"})
		return
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Entity enhancement completed in %.2f seconds. Processed: %d, Enhanced: %d, Removed: %d\n",
		elapsed.Seconds(), result.ProcessedCount, len(result.EnhancedEntities), result.RemovedCount)

	c.JSON(http.StatusOK, result)
}

// ServiceInfo handles the service information endpoint
func (h *Handler) ServiceInfo(c *gin.Context) {
	response := models.ServiceInfo{
		Title:       "LLM Inference Service",
		Description: "A service for inference tasks using large language models.",
		Version:     "0.1.0",
		Endpoints: map[string]string{
			"similarity":       "/api/v1/similarity",
			"enhance_entities": "/api/v1/enhance-entities",
			"health":           "/health",
		},
	}
	c.JSON(http.StatusOK, response)
}
