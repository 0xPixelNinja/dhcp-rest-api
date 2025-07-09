package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns service status - no auth required
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "dhcp-rest-api",
		"message": "Service is running",
	})
}
