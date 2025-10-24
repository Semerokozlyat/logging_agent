package logaggregator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/grafana/loki/v3/clients/pkg/promtail/api"
	lokiClient "github.com/grafana/loki/v3/clients/pkg/promtail/client"
	"github.com/grafana/loki/v3/pkg/logproto"
	"github.com/prometheus/common/model"

	"github.com/Semerokozlyat/logging_agent/internal/config"
)

const (
	backendTypeStdout = "stdout"
	backendTypeLoki   = "loki"

	appNameLabel = "logging-agent"
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
	logChan     <-chan LogEntry
	backendType string
	lokiClient  lokiClient.Client
}

func New(cfg *config.Config, logChan <-chan LogEntry) (*Aggregator, error) {
	if logChan == nil {
		return nil, errors.New("log channel is not initialized")
	}

	a := Aggregator{
		logChan:     logChan,
		backendType: backendTypeStdout,
	}

	if cfg.Loki.URL.String() != "" {
		var err error
		a.lokiClient, err = NewLokiClient(cfg.Loki, cfg.Agent.Collection)
		if err != nil {
			return nil, fmt.Errorf("init Loki client: %w", err)
		}
		a.backendType = backendTypeLoki
	}
	return &a, nil
}

func (a *Aggregator) processLogEntry(entry LogEntry) error {
	if a.backendType == backendTypeLoki {
		// Convert LogEntry to api.Entry
		labels := model.LabelSet{
			"app":    model.LabelValue(appNameLabel),
			"node":   model.LabelValue(entry.NodeName),
			"source": model.LabelValue(entry.Source),
			"level":  model.LabelValue(entry.Level),
		}

		lokiEntry := api.Entry{
			Labels: labels,
			Entry: logproto.Entry{
				Timestamp: entry.Timestamp,
				Line:      entry.Message,
			},
		}

		a.lokiClient.Chan() <- lokiEntry
		return nil
	}
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
		case entry, ok := <-a.logChan:
			if !ok {
				log.Print("log aggregator channel is closed, exiting")
				return
			}
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

func (a *Aggregator) Stop() {
	if a.lokiClient != nil {
		a.lokiClient.Stop()
	}
}
