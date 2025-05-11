package config

// GetDevelopmentConfig returns development environment settings
func GetDevelopmentConfig() *Config {
	return &Config{
		Environment: "development",
		DataDir:     "./api/data",
		Port:        "8080",
		BaseURL:     "http://localhost:8080",
	}
}
