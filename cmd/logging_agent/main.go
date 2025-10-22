package main

import (
	"fmt"
	"log"

	"github.com/Semerokozlyat/logging_agent/pkg/agent"
)

func main() {
	log.Println("Logging Agent Starting...")

	a := agent.New()
	if err := a.Run(); err != nil {
		log.Fatalf("Failed to run agent: %v", err)
	}

	fmt.Println("Logging Agent Running")
}
