package handlers

import (
	"llm-inference/models"

	"github.com/gofiber/fiber/v2"
)

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
