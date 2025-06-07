package main

import (
	"log"
	"net/http"
	"os"

	"github.com/0xPixelNinja/dhcp-rest-api/auth"
	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/0xPixelNinja/dhcp-rest-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration
	config.LoadConfig()

	// Set Gin mode (e.g., release, debug, test)
	// Default is debug mode. For production, set to release mode.
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production

	// Initialize Gin router
	router := gin.Default()

	// Simple health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Apply authentication middleware to protected route groups
	authedRoutes := router.Group("/") // Apply to root or specific groups as needed
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

	// Start server
	port := os.Getenv("PORT") // Allow port to be set from environment
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
