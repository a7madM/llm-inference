package handlers

import (
	"fmt"
	"time"

	"llm-inference/models"

	"github.com/gofiber/fiber/v2"
)

// Verifier handles the entity verification endpoint
func (h *Handler) Verifier(c *fiber.Ctx) error {
	var input models.EntityVerifierRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request body"})
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
	fmt.Printf("Entity enhancement completed in %.2f seconds, and results -> %v\n", elapsed.Seconds(), result)

	return c.Status(fiber.StatusOK).JSON(result)
}
