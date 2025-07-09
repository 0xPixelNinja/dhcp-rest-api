package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck provides a simple health check endpoint
// This endpoint does not require authentication
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "dhcp-rest-api",
		"message": "Service is running",
	})
}
