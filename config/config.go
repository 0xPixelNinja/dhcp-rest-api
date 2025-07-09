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
	Environment        string
	Port               string
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig loads configuration from environment variables
// It looks for a .env file first, then environment variables.
func LoadConfig() {
	// Only load .env file in development mode
	env := getEnv("ENVIRONMENT", "development")
	if env == "development" {
		_ = godotenv.Load()
	}

	AppConfig.Environment = env
	AppConfig.DhcpConfPath = getEnv("DHCP_CONF_PATH", "/etc/dhcp/dhcpd.conf")
	AppConfig.InterfacesConfPath = getEnv("INTERFACES_CONF_PATH", "/etc/default/isc-dhcp-server")
	AppConfig.TokenFilePath = getEnv("TOKEN_FILE_PATH", "/etc/dhcp-rest-api/token")
	AppConfig.Port = getEnv("PORT", "8080")

	// Load token from file if exists, otherwise use env var
	if tokenFromFile := loadTokenFromFile(); tokenFromFile != "" {
		AppConfig.TokenSecret = tokenFromFile
		if AppConfig.Environment == "development" {
			log.Println("Token loaded from file")
		}
	} else {
		AppConfig.TokenSecret = getEnv("TOKEN_SECRET", "")
		if AppConfig.TokenSecret == "" {
			log.Fatal("TOKEN_SECRET must be set in production")
		}
		if AppConfig.Environment == "development" {
			log.Println("Token loaded from environment variable")
		}
	}

	// Validate required configuration in production
	if AppConfig.Environment == "production" {
		validateProductionConfig()
	}

	if AppConfig.Environment == "development" {
		log.Println("Configuration loaded successfully")
		log.Printf("Environment: %s", AppConfig.Environment)
		log.Printf("DHCP_CONF_PATH: %s", AppConfig.DhcpConfPath)
		log.Printf("INTERFACES_CONF_PATH: %s", AppConfig.InterfacesConfPath)
		log.Printf("TOKEN_FILE_PATH: %s", AppConfig.TokenFilePath)
		log.Printf("PORT: %s", AppConfig.Port)
	}
}

// validateProductionConfig ensures all required settings are present in production
func validateProductionConfig() {
	if AppConfig.TokenSecret == "" {
		log.Fatal("TOKEN_SECRET is required in production")
	}

	// Check if config paths exist and are accessible
	if _, err := os.Stat(AppConfig.DhcpConfPath); os.IsNotExist(err) {
		log.Printf("Warning: DHCP config path does not exist: %s", AppConfig.DhcpConfPath)
	}

	if _, err := os.Stat(AppConfig.InterfacesConfPath); os.IsNotExist(err) {
		log.Printf("Warning: Interfaces config path does not exist: %s", AppConfig.InterfacesConfPath)
	}

	log.Printf("Production configuration validated successfully")
}

// IsProduction returns true if running in production mode
func IsProduction() bool {
	return AppConfig.Environment == "production"
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
		if !IsProduction() {
			log.Println("Token updated successfully")
		}
	}
	return err
}
