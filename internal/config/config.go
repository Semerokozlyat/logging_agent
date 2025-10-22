package config

import (
	"os"
	"time"
)

// AppConfig represents the application configuration
type Config struct {
	Agent      Agent      `json:"agent" yaml:"agent"`
	HTTPServer HTTPServer `json:"http_server" yaml:"httpServer"`
}

type Agent struct {
	LogLevel           string
	OutputPath         string
	LogPaths           []string
	CollectionInterval time.Duration
	BatchSize          int
	MaxLineLength      int
	NodeName           string
	PodName            string
	Namespace          string
}

type HTTPServer struct {
	Address      string        `json:"address" yaml:"address"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idleTimeout"`
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
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
