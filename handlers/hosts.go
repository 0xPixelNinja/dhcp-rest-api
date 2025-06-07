package handlers

import (
	"net/http"

	"log"

	"github.com/0xPixelNinja/dhcp-rest-api/models"
	"github.com/0xPixelNinja/dhcp-rest-api/services"
	"github.com/gin-gonic/gin"
)

// ListHosts retrieves a list of all DHCP hosts.
// Corresponds to GET /hosts/
func ListHosts(c *gin.Context) {
	hosts, err := services.ListHosts()
	if err != nil {
		log.Printf("Error listing hosts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list hosts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hosts": hosts})
}

// AddHost adds a new DHCP host.
// Corresponds to POST /hosts/
func AddHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := services.AddHost(host); err != nil {
		log.Printf("Error adding host %s: %v", host.Name, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add host"}) // Python app returns 400
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host added successfully"})
}

// UpdateHost updates an existing DHCP host.
// Corresponds to PUT /hosts/{name}
func UpdateHost(c *gin.Context) {
	hostName := c.Param("name")
	var hostUpdate models.HostUpdate

	if err := c.ShouldBindJSON(&hostUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := services.UpdateHost(hostName, hostUpdate); err != nil {
		log.Printf("Error updating host %s: %v", hostName, err)
		// Check if the error is because the host was not found, similar to Python raising HTTPException 400
		// For now, defaulting to a generic 400 as the Python app does.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update host"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host updated successfully"})
}

// DeleteHost deletes an existing DHCP host.
// Corresponds to DELETE /hosts/{name}
func DeleteHost(c *gin.Context) {
	hostName := c.Param("name")

	if err := services.DeleteHost(hostName); err != nil {
		log.Printf("Error deleting host %s: %v", hostName, err)
		// Python app returns 400 if delete fails for reasons other than not found (where it succeeds)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete host"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host deleted successfully"})
}
