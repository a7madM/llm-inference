package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"llm-inference/models"
	"llm-inference/services"

	"github.com/gin-gonic/gin"
)

// Handler holds the service dependencies
type Handler struct {
	nerService        *services.NERService
	similarityService *services.SimilarityService
}

func NewHandler(nerService *services.NERService, similarityService *services.SimilarityService) *Handler {
	return &Handler{
		nerService:        nerService,
		similarityService: similarityService,
	}
}

// ExtractEntities handles the entities extraction endpoint
func (h *Handler) ExtractEntities(c *gin.Context) {
	var input models.InputText
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Log request
	textPreview := input.Text
	if len(textPreview) > 50 {
		textPreview = textPreview[:50] + "..."
	}
	fmt.Printf("Received text for NER: %s\n", textPreview)

	startTime := time.Now()

	// Extract entities
	entities, err := h.nerService.ExtractEntities(input.Text)
	if err != nil {
		log.Printf("Error extracting entities: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Unable to extract entities"})
		return
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Extraction completed in %.2f seconds.\n", elapsed.Seconds())

	c.JSON(http.StatusOK, entities)
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	response := models.HealthResponse{
		Status:  "healthy",
		Service: "Local Multilingual NER Service",
		Uptime:  time.Now().UTC().Format(time.RFC3339),
		Version: "1.2.0",
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

// ServiceInfo handles the service information endpoint
func (h *Handler) ServiceInfo(c *gin.Context) {
	response := models.ServiceInfo{
		Title:       "Local Multilingual NER Service",
		Description: "Extracts Persons, Locations, Organizations, Events from text (Arabic, English, German). Preserves original language.",
		Version:     "1.2.0",
		Endpoints: map[string]string{
			"entities":  "/api/v1/entities",
			"sentiment": "/api/v1/similarity",
			"health":    "/health",
		},
	}
	c.JSON(http.StatusOK, response)
}
