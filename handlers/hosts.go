package handlers

import (
	"log"
	"net/http"

	"github.com/0xPixelNinja/dhcp-rest-api/models"
	"github.com/0xPixelNinja/dhcp-rest-api/services"
	"github.com/gin-gonic/gin"
)

func ListHosts(c *gin.Context) {
	hosts, err := services.ListHosts()
	if err != nil {
		log.Printf("Error listing hosts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list hosts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hosts": hosts})
}

func AddHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := services.AddHost(host); err != nil {
		log.Printf("Error adding host %s: %v", host.Name, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add host"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host added successfully"})
}

func UpdateHost(c *gin.Context) {
	hostName := c.Param("name")
	var hostUpdate models.HostUpdate

	if err := c.ShouldBindJSON(&hostUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := services.UpdateHost(hostName, hostUpdate); err != nil {
		log.Printf("Error updating host %s: %v", hostName, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update host"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host updated successfully"})
}

func DeleteHost(c *gin.Context) {
	hostName := c.Param("name")

	if err := services.DeleteHost(hostName); err != nil {
		log.Printf("Error deleting host %s: %v", hostName, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete host"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Host deleted successfully"})
}
