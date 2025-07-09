package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DhcpConfPath       string
	InterfacesConfPath string
	TokenSecret        string
	TokenFilePath      string
	Environment        string
	Port               string
}

var AppConfig Config

func LoadConfig() {
	// Only load .env file in development
	env := getEnv("ENVIRONMENT", "development")
	if env == "development" {
		_ = godotenv.Load()
	}

	AppConfig.Environment = env
	AppConfig.DhcpConfPath = getEnv("DHCP_CONF_PATH", "/etc/dhcp/dhcpd.conf")
	AppConfig.InterfacesConfPath = getEnv("INTERFACES_CONF_PATH", "/etc/default/isc-dhcp-server")
	AppConfig.TokenFilePath = getEnv("TOKEN_FILE_PATH", "/etc/dhcp-rest-api/token")
	AppConfig.Port = getEnv("PORT", "8080")

	// Try to load token from file first, then environment, then auto-generate
	if tokenFromFile := loadTokenFromFile(); tokenFromFile != "" {
		AppConfig.TokenSecret = tokenFromFile
		if AppConfig.Environment == "development" {
			log.Println("Token loaded from file")
		}
	} else if envToken := getEnv("TOKEN_SECRET", ""); envToken != "" {
		AppConfig.TokenSecret = envToken
		if AppConfig.Environment == "development" {
			log.Println("Token loaded from environment variable")
		}
	} else {
		// Auto-generate a secure token for first-time setup
		generatedToken, err := generateSecureToken()
		if err != nil {
			log.Fatalf("Failed to generate secure token: %v", err)
		}

		AppConfig.TokenSecret = generatedToken

		// Save the generated token to file
		if err := SaveToken(generatedToken); err != nil {
			log.Printf("Warning: Failed to save auto-generated token to file: %v", err)
			log.Printf("Auto-generated token (save this!): %s", generatedToken)
		} else {
			log.Printf("Auto-generated secure token and saved to: %s", AppConfig.TokenFilePath)
			log.Printf("Your API token: %s", generatedToken)
		}

		log.Println("IMPORTANT: Save this token! You'll need it to access the API.")
		log.Println("You can also retrieve it later from the token file or change it via the API.")
	}

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

func validateProductionConfig() {
	if AppConfig.TokenSecret == "" {
		log.Fatal("TOKEN_SECRET is required in production")
	}

	// Check if config files exist
	if _, err := os.Stat(AppConfig.DhcpConfPath); os.IsNotExist(err) {
		log.Printf("Warning: DHCP config path does not exist: %s", AppConfig.DhcpConfPath)
	}

	if _, err := os.Stat(AppConfig.InterfacesConfPath); os.IsNotExist(err) {
		log.Printf("Warning: Interfaces config path does not exist: %s", AppConfig.InterfacesConfPath)
	}

	log.Printf("Production configuration validated successfully")
}

// generateSecureToken creates a cryptographically secure random token
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsProduction() bool {
	return AppConfig.Environment == "production"
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func loadTokenFromFile() string {
	if content, err := os.ReadFile(AppConfig.TokenFilePath); err == nil {
		return strings.TrimSpace(string(content))
	}
	return ""
}

func SaveToken(token string) error {
	dir := filepath.Dir(AppConfig.TokenFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	err := os.WriteFile(AppConfig.TokenFilePath, []byte(token), 0600)
	if err == nil {
		AppConfig.TokenSecret = token
		if !IsProduction() {
			log.Println("Token updated successfully")
		}
	}
	return err
}
