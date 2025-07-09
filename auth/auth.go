package auth

import (
	"net/http"
	"strings"

	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for valid Bearer token in Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid authorization header format. Expected Bearer token."})
			c.Abort()
			return
		}

		token := parts[1]
		if token != config.AppConfig.TokenSecret {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or missing token."})
			c.Abort()
			return
		}

		c.Next()
	}
}
