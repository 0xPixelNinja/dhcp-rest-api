package main

import (
	"log"

	"github.com/0xPixelNinja/dhcp-rest-api/auth"
	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/0xPixelNinja/dhcp-rest-api/handlers"
	"github.com/0xPixelNinja/dhcp-rest-api/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router with production settings
	r := gin.New()

	// Add recovery middleware to handle panics gracefully
	r.Use(gin.Recovery())

	// Add custom logging middleware for production
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"}, // Skip logging health check requests to reduce noise
	}))

	// Add security headers middleware
	r.Use(middleware.SecurityHeaders())

	// Add CORS middleware
	r.Use(middleware.CORS())

	// Create rate limiter
	rateLimiter := middleware.DefaultRateLimiter()

	// Apply rate limiting to all routes
	r.Use(rateLimiter.Middleware())

	// Public health check endpoint (with rate limiting but no authentication)
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
	log.Printf("Starting DHCP REST API server on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
