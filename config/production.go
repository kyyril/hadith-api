package config

import "os"

// GetProductionConfig returns production environment settings
func GetProductionConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Environment: "production",
		DataDir:     "./data",
		Port:        port,
		BaseURL:     os.Getenv("https://go-hadist.vercel.app"),
	}
}
