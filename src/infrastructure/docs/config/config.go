package config

import (
	"os"
)

// Config contains all the configuration of the Swagger server
type Config struct {
	// Swagger server
	SwaggerPort string

	// API information
	APITitle       string
	APIVersion     string
	APIDescription string
	APIHost        string
	APIBasePath    string
}

// Load the configuration from environment variables
func Load() *Config {
	return &Config{
		SwaggerPort:    getEnv("SWAGGER_PORT", "8081"),
		APITitle:       getEnv("API_TITLE", "GoProjectSkeleton API"),
		APIVersion:     getEnv("API_VERSION", "1.0"),
		APIDescription: getEnv("API_DESCRIPTION", "API documentation for GoProjectSkeleton"),
		APIHost:        getEnv("API_HOST", "localhost:8080"),
		APIBasePath:    getEnv("API_BASE_PATH", ""),
	}
}

// getEnv get an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
