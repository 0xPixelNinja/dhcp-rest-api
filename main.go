package main

import (
	"log"
	"os"

	"github.com/0xPixelNinja/dhcp-rest-api/auth"
	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/0xPixelNinja/dhcp-rest-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Create Gin router
	r := gin.Default()

	// Public health check endpoint (no authentication required)
	r.GET("/health", handlers.HealthCheck)

	// Apply authentication middleware to all protected routes
	authedRoutes := r.Group("/")
	authedRoutes.Use(auth.AuthMiddleware())

	// Host routes
	hostRoutes := authedRoutes.Group("/hosts")
	{
		hostRoutes.GET("/", handlers.ListHosts)
		hostRoutes.POST("/", handlers.AddHost)
		hostRoutes.PUT("/:name", handlers.UpdateHost)
		hostRoutes.DELETE("/:name", handlers.DeleteHost)
	}

	// Interface routes
	interfaceRoutes := authedRoutes.Group("/interfaces")
	{
		interfaceRoutes.GET("/", handlers.ListInterfaces)
		interfaceRoutes.POST("/", handlers.AddInterface)      // Body contains type and interface
		interfaceRoutes.DELETE("/", handlers.DeleteInterface) // Body contains type and interface
	}

	// Token management route
	authedRoutes.PUT("/token", handlers.UpdateToken)

	// Start server
	port := os.Getenv("PORT") // Allow port to be set from environment
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
