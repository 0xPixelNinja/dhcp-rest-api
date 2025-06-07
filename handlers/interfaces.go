package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/0xPixelNinja/dhcp-rest-api/models"
	"github.com/0xPixelNinja/dhcp-rest-api/services"
	"github.com/gin-gonic/gin"
)

// ListInterfaces retrieves the current list of network interfaces from the config.
// Corresponds to GET /interfaces/
func ListInterfaces(c *gin.Context) {
	interfaces, err := services.GetInterfaces()
	if err != nil {
		log.Printf("Error listing interfaces: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list interfaces"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"interfaces": interfaces})
}

// AddInterface adds a network interface to the configuration.
// Corresponds to POST /interfaces/
func AddInterface(c *gin.Context) {
	var op models.InterfaceOperation
	if err := c.ShouldBindJSON(&op); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := services.AddInterface(op.Type, op.Interface); err != nil {
		log.Printf("Error adding interface %s (type %s): %v", op.Interface, op.Type, err)
		// Python app returns 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add interface."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Interface %s added to INTERFACES%s successfully.", op.Interface, op.Type)})
}

// DeleteInterface removes a network interface from the configuration.
// Corresponds to DELETE /interfaces/
func DeleteInterface(c *gin.Context) {
	var op models.InterfaceOperation
	// For DELETE requests with a body, ShouldBindJSON is appropriate.
	// If the parameters were in the query string or path, we'd use c.Query() or c.Param().
	if err := c.ShouldBindJSON(&op); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := services.DeleteInterface(op.Type, op.Interface); err != nil {
		log.Printf("Error deleting interface %s (type %s): %v", op.Interface, op.Type, err)
		// Python app returns 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete interface."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Interface %s removed from INTERFACES%s successfully.", op.Interface, op.Type)})
}
