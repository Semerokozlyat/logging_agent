package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Semerokozlyat/logging_agent/internal/agent"
	"github.com/Semerokozlyat/logging_agent/internal/config"
)

func main() {
	log.Println("Logging Agent is Starting...")

	var configPath = flag.String("config", "/etc/logging-agent/config.yaml", "Path to config file.")
	flag.Parse()

	appConfig, err := config.New(*configPath)
	if err != nil {
		log.Fatal("failed to initialize app config", err)
	}

	// Create agent
	agent := agent.New(appConfig)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	// Start agent in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := agent.Run(ctx); err != nil {
			errChan <- err
		}
	}()

	fmt.Println("Logging Agent is Running")

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		cancel()
	case err := <-errChan:
		log.Fatalf("Agent error: %v", err)
	}

	log.Println("Logging Agent Stopped")
}
