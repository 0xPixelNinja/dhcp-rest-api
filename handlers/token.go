package handlers

import (
	"log"
	"net/http"

	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/gin-gonic/gin"
)

type TokenUpdateRequest struct {
	Token string `json:"token" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// GetToken returns the current authentication token (requires auth)
func GetToken(c *gin.Context) {
	c.JSON(http.StatusOK, TokenResponse{
		Token: config.AppConfig.TokenSecret,
	})
}

// UpdateToken allows updating the authentication token
func UpdateToken(c *gin.Context) {
	var req TokenUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token cannot be empty"})
		return
	}

	if err := config.SaveToken(req.Token); err != nil {
		log.Printf("Error saving new token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	log.Printf("Authentication token updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Token updated successfully"})
}
