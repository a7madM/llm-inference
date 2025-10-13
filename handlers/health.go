package handlers

import (
	"time"

	"llm-inference/models"

	"github.com/gofiber/fiber/v2"
)

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
