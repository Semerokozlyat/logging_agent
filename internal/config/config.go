package config

import (
	"os"
)

// AppConfig represents the application configuration
type AppConfig struct {
	Port     string
	LogLevel string
	Debug    bool
}

// Load loads configuration from environment variables
func Load() *AppConfig {
	return &AppConfig{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Debug:    getEnv("DEBUG", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
