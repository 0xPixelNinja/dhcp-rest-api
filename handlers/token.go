package handlers

import (
	"log"
	"net/http"

	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/gin-gonic/gin"
)

// TokenUpdateRequest represents the request payload for token updates
type TokenUpdateRequest struct {
	Token string `json:"token" binding:"required"`
}

// UpdateToken allows updating the authentication token
// This endpoint uses the current token for authentication
func UpdateToken(c *gin.Context) {
	var req TokenUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Validate token is not empty
	if req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token cannot be empty"})
		return
	}

	// Save the new token
	if err := config.SaveToken(req.Token); err != nil {
		log.Printf("Error saving new token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	log.Printf("Authentication token updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Token updated successfully"})
}
