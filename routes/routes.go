package routes

import (
	"net/http"
	"os"

	"llm-inference/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main router
func SetupRouter(handler *handlers.Handler) *gin.Engine {
	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	// API routes
	setupAPIRoutes(router, handler)

	// Health and info routes
	setupHealthRoutes(router, handler)

	return router
}

// setupAPIRoutes configures the API v1 routes
func setupAPIRoutes(router *gin.Engine, handler *handlers.Handler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/entities", handler.ExtractEntities)
		v1.POST("/sentiment", handler.AnalyzeSentiment)
	}
}

// setupHealthRoutes configures health check and info routes
func setupHealthRoutes(router *gin.Engine, handler *handlers.Handler) {
	router.GET("/health", handler.HealthCheck)
	router.GET("/", handler.ServiceInfo)
}

// corsMiddleware provides CORS support
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
