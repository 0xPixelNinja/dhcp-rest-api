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
	config.LoadConfig()

	// Use production mode - no debug output
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Basic middleware stack
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"}, // health checks spam the logs
	}))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())

	// Rate limiting for all endpoints
	rateLimiter := middleware.DefaultRateLimiter()
	r.Use(rateLimiter.Middleware())

	// Public endpoint - no auth needed
	r.GET("/health", handlers.HealthCheck)

	// Everything else requires authentication
	authedRoutes := r.Group("/")
	authedRoutes.Use(auth.AuthMiddleware())

	// DHCP host management
	hostRoutes := authedRoutes.Group("/hosts")
	{
		hostRoutes.GET("/", handlers.ListHosts)
		hostRoutes.POST("/", handlers.AddHost)
		hostRoutes.PUT("/:name", handlers.UpdateHost)
		hostRoutes.DELETE("/:name", handlers.DeleteHost)
	}

	// Network interface management
	interfaceRoutes := authedRoutes.Group("/interfaces")
	{
		interfaceRoutes.GET("/", handlers.ListInterfaces)
		interfaceRoutes.POST("/", handlers.AddInterface)
		interfaceRoutes.DELETE("/", handlers.DeleteInterface)
	}

	// Token management
	// authedRoutes.GET("/token", handlers.GetToken)
	// authedRoutes.PUT("/token", handlers.UpdateToken)

	log.Printf("Starting DHCP REST API server on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
