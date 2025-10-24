package config

import (
	"fmt"
	"os"
	"time"

	"go.yaml.in/yaml/v2"
)

const (
	defaultLogLevel = "info"

	nodeNameEnvVar  = "NODE_NAME"
	podNameEnvVar   = "POD_NAME"
	namespaceEnvVar = "POD_NAMESPACE"
	logLevelEnvVar  = "LOG_LEVEL"
)

// AppConfig represents the application configuration
type Config struct {
	Agent      Agent      `json:"agent" yaml:"agent"`
	HTTPServer HTTPServer `json:"http_server" yaml:"httpServer"`
}

type Agent struct {
	LogLevel   string     `json:"log_level" yaml:"logLevel"`
	OutputPath string     `json:"output_path" yaml:"outputPath"`
	Collection Collection `json:"collection" yaml:"collection"`
	NodeName   string     `json:"node_name" yaml:"nodeName"`
	PodName    string     `json:"pod_name" yaml:"podName"`
	Namespace  string     `json:"namespace" yaml:"namespace"`
}

type Collection struct {
	LogChanSize   int           `json:"log_chan_size" yaml:"logChanSize"`
	LogPaths      []string      `json:"log_paths" yaml:"logPaths"`
	Interval      time.Duration `json:"interval" yaml:"interval"`
	BatchSize     int           `json:"batch_size" yaml:"batchSize"`
	MaxLineLength int           `json:"max_line_length" yaml:"maxLineLength"`
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
	cfg.Agent.NodeName = getEnv(nodeNameEnvVar, cfg.Agent.NodeName)
	cfg.Agent.PodName = getEnv(podNameEnvVar, cfg.Agent.PodName)
	cfg.Agent.Namespace = getEnv(namespaceEnvVar, cfg.Agent.Namespace)
	cfg.Agent.LogLevel = getEnv(logLevelEnvVar, defaultLogLevel)

	return &cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
