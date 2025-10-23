package config

import (
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v2"
)

const (
	defaultLogLevel = "info"
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

// New loads configuration from file and returns config instance.
func New(configFilePath string) (*Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", configFilePath, err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file %q data: %w", configFilePath, err)
	}

	// Redefine config values from environment vars
	cfg.Agent.NodeName = getEnv("NODE_NAME", cfg.Agent.NodeName)
	cfg.Agent.PodName = getEnv("POD_NAME", cfg.Agent.PodName)
	cfg.Agent.Namespace = getEnv("POD_NAMESPACE", cfg.Agent.Namespace)
	cfg.Agent.LogLevel = getEnv("LOG_LEVEL", defaultLogLevel)

	return &cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
