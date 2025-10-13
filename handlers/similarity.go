package handlers

import (
	"fmt"

	"llm-inference/models"

	"github.com/gofiber/fiber/v2"
)

// ComputeSimilarity handles the similarity endpoint
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
