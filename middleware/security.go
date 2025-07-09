package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers for production
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Remove server information
		c.Header("Server", "")

		// Add Content Security Policy for API
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")

		// Prevent caching of sensitive endpoints
		if c.Request.URL.Path != "/health" {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// CORS configures Cross-Origin Resource Sharing for production
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// In production, you should whitelist specific origins
		// For now, we'll be restrictive and only allow specific origins
		allowedOrigins := []string{
			"*",
		}

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
