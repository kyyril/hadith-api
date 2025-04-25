package config

import "os"

// Config holds environment-specific configuration
type Config struct {
	Environment string
	DataDir     string
	Port        string
	BaseURL     string
	// Add other config fields as needed
}

// GetConfig returns environment-specific configuration
func GetConfig() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // Default environment
	}

	if env == "production" {
		return GetProductionConfig()
	}
	return GetDevelopmentConfig()
}
