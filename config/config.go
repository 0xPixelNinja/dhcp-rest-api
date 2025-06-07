package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DhcpConfPath       string
	InterfacesConfPath string
	TokenSecret        string
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig loads configuration from environment variables
// It looks for a .env file first, then environment variables.
func LoadConfig() {
	// Attempt to load .env file, but don't fail if it's not present
	// as environment variables might be set directly (e.g., in Docker)
	_ = godotenv.Load()

	AppConfig.DhcpConfPath = getEnv("DHCP_CONF_PATH", "/etc/dhcp/dhcpd.conf")
	AppConfig.InterfacesConfPath = getEnv("INTERFACES_CONF_PATH", "/etc/default/isc-dhcp-server")
	AppConfig.TokenSecret = getEnv("TOKEN_SECRET", "your-secret-token") // Default matches Python app

	log.Println("Configuration loaded successfully")
	log.Printf("DHCP_CONF_PATH: %s", AppConfig.DhcpConfPath)
	log.Printf("INTERFACES_CONF_PATH: %s", AppConfig.InterfacesConfPath)
	// Avoid logging the token secret for security reasons
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
