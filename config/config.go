package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DhcpConfPath       string
	InterfacesConfPath string
	TokenSecret        string
	TokenFilePath      string
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
	AppConfig.TokenFilePath = getEnv("TOKEN_FILE_PATH", "/etc/dhcp-rest-api/token")

	// Load token from file if exists, otherwise use env var
	if tokenFromFile := loadTokenFromFile(); tokenFromFile != "" {
		AppConfig.TokenSecret = tokenFromFile
		log.Println("Token loaded from file")
	} else {
		AppConfig.TokenSecret = getEnv("TOKEN_SECRET", "your-secret-token")
		log.Println("Token loaded from environment variable")
	}

	log.Println("Configuration loaded successfully")
	log.Printf("DHCP_CONF_PATH: %s", AppConfig.DhcpConfPath)
	log.Printf("INTERFACES_CONF_PATH: %s", AppConfig.InterfacesConfPath)
	log.Printf("TOKEN_FILE_PATH: %s", AppConfig.TokenFilePath)
	// Avoid logging the token secret for security reasons
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// loadTokenFromFile loads the authentication token from file
func loadTokenFromFile() string {
	if content, err := ioutil.ReadFile(AppConfig.TokenFilePath); err == nil {
		return strings.TrimSpace(string(content))
	}
	return ""
}

// SaveToken saves a new authentication token to file and updates in-memory config
func SaveToken(token string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(AppConfig.TokenFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write token to file with restricted permissions
	err := ioutil.WriteFile(AppConfig.TokenFilePath, []byte(token), 0600)
	if err == nil {
		// Update in-memory configuration
		AppConfig.TokenSecret = token
		log.Println("Token updated successfully")
	}
	return err
}
