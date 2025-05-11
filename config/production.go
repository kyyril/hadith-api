// Modified config/production.go
package config

import (
	"os"
	"path/filepath"
)

// GetProductionConfig returns production environment settings
func GetProductionConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// For Vercel serverless functions, we need to use absolute paths
	// __dirname equivalent in Go
	execDir, _ := os.Getwd()
	dataDir := filepath.Join(execDir, "data")

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://hadith-api-go.vercel.app"
	}

	return &Config{
		Environment: "production",
		DataDir:     dataDir,
		Port:        port,
		BaseURL:     baseURL,
	}
}
