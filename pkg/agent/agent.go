package agent

import (
	"fmt"
)

// Agent represents the logging agent
type Agent struct {
	config Config
}

// Config holds agent configuration
type Config struct {
	LogLevel   string
	OutputPath string
}

// New creates a new Agent instance
func New() *Agent {
	return &Agent{
		config: Config{
			LogLevel:   "info",
			OutputPath: "/var/log/agent.log",
		},
	}
}

// Run starts the agent
func (a *Agent) Run() error {
	fmt.Printf("Agent running with config: %+v\n", a.config)
	// Add your agent logic here
	return nil
}
