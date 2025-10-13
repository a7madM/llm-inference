package routes

import (
	"llm-inference/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupRouter configures and returns the main Fiber app
func SetupRouter(handler *handlers.Handler) *fiber.App {
	app := fiber.New()

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	// API routes
	v1 := app.Group("/api/v1")
	v1.Get("/similarity", handler.ComputeSimilarity)
	v1.Post("/verify", handler.Verifier)

	// Health and info routes
	app.Get("/health", handler.HealthCheck)
	app.Get("/", handler.ServiceInfo)

	return app
}
