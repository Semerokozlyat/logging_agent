package logaggregator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	NodeName  string    `json:"node_name"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

type Aggregator struct {
	logChan <-chan LogEntry
}

func New(logChan <-chan LogEntry) (*Aggregator, error) {
	if logChan == nil {
		return nil, errors.New("log channel is not initialized")
	}
	return &Aggregator{
		logChan: logChan,
	}, nil
}

func (a *Aggregator) processLogEntry(entry LogEntry) error {
	fmt.Printf("[%s] [%s] [%s] %s: %s\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.NodeName,
		entry.Source,
		entry.Message,
	)
	return nil
}

func (a *Aggregator) Run(ctx context.Context) {
	for {
		select {
		case entry := <-a.logChan:
			err := a.processLogEntry(entry)
			if err != nil {
				log.Printf("failed to process log entry: %s", err)
			}
		case <-ctx.Done():
			log.Print("aggregator caught stop signal from context")
			return
		}
	}
}
