package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Semerokozlyat/logging_agent/pkg/agent"
)

func main() {
	log.Println("Logging Agent Starting...")

	// Create agent
	a := agent.New()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start agent in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := a.Run(); err != nil {
			errChan <- err
		}
	}()

	fmt.Println("Logging Agent Running - Press Ctrl+C to stop")

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		a.Stop()
	case err := <-errChan:
		log.Fatalf("Agent error: %v", err)
	}

	log.Println("Logging Agent Stopped")
}
