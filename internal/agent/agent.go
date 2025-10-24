package agent

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Semerokozlyat/logging_agent/internal/config"
	"github.com/Semerokozlyat/logging_agent/internal/httpserver"
	"github.com/Semerokozlyat/logging_agent/internal/pkg/logaggregator"
	"github.com/Semerokozlyat/logging_agent/internal/pkg/metrics"
)

type Meta struct {
	NodeName  string
	PodName   string
	Namespace string
}

type LogAggregator interface {
	Run(ctx context.Context)
	Stop()
}

// Agent represents the logging agent
type Agent struct {
	logFiles map[string]*os.File
	logMutex sync.RWMutex

	healthCheckServer *http.Server

	logChan       chan logaggregator.LogEntry
	logAggregator LogAggregator

	OutputPath         string
	LogPaths           []string
	CollectionInterval time.Duration
	BatchSize          int
	MaxLineLength      int
	Meta               Meta
}

// New creates a new Agent instance
func New(cfg *config.Config) (*Agent, error) {

	// Init metrics
	metrics.InitMetricsCollector()

	log.Printf("Initializing logging agent with configuration: %+v", cfg)

	logChan := make(chan logaggregator.LogEntry, cfg.Agent.Collection.LogChanSize)
	logAggregator, err := logaggregator.New(cfg, logChan)
	if err != nil {
		return nil, fmt.Errorf("init log aggregator: %w", err)
	}

	return &Agent{
		OutputPath:         cfg.Agent.OutputPath,
		LogPaths:           cfg.Agent.Collection.LogPaths,
		CollectionInterval: cfg.Agent.Collection.Interval,
		BatchSize:          cfg.Agent.Collection.BatchSize,
		MaxLineLength:      cfg.Agent.Collection.MaxLineLength,
		Meta: Meta{
			NodeName:  cfg.Agent.NodeName,
			PodName:   cfg.Agent.PodName,
			Namespace: cfg.Agent.Namespace,
		},
		logFiles: make(map[string]*os.File),

		logChan:       logChan,
		logAggregator: logAggregator,

		healthCheckServer: httpserver.NewHealthCheckServer(&cfg.HTTPServer),
	}, nil
}

// Run starts the agent
func (a *Agent) Run(ctx context.Context) error {
	log.Printf("Starting Logging Agent on node: %s", a.Meta.NodeName)

	var wg sync.WaitGroup

	// Start health check server
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if p := recover(); p != nil {
				log.Fatal(fmt.Sprintf("Panic in healthcheck HTTP server goroutine: %s %s", p, debug.Stack()))
			}
		}()
		log.Printf("Starting healthcheck HTTP server on %s", a.healthCheckServer.Addr)
		err := a.healthCheckServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(fmt.Sprintf("Failed to start healthcheck HTTP server on %s: %s", a.healthCheckServer.Addr, err))
		}
	}()

	// Start log aggregator
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if p := recover(); p != nil {
				log.Fatal(fmt.Sprintf("Panic in logs aggregation goroutine: %s %s", p, debug.Stack()))
			}
		}()
		log.Printf("Start log aggregator")
		a.logAggregator.Run(ctx)
	}()

	// Start log collection
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if p := recover(); p != nil {
				log.Fatal(fmt.Sprintf("Panic in logs collection goroutine: %s %s", p, debug.Stack()))
			}
		}()
		log.Printf("Start logs collection")
		a.collectLogs(ctx)
	}()

	// Wait for context cancellation
	<-ctx.Done()

	log.Println("Shutting down agent...")
	wg.Wait()

	// Stop the agent
	a.Stop()

	log.Println("Agent stopped gracefully")
	return nil
}

// Stop gracefully stops the agent
func (a *Agent) Stop() {
	log.Println("Stop signal received")
	a.closeLogFiles()
	a.logAggregator.Stop()
}

// collectLogs collects logs from specified paths
func (a *Agent) collectLogs(ctx context.Context) {
	ticker := time.NewTicker(a.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.scanAndProcessLogs()
		}
	}
}

// scanAndProcessLogs scans log paths and processes new log entries
func (a *Agent) scanAndProcessLogs() {
	for _, pattern := range a.LogPaths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("Error globbing pattern %s: %v", pattern, err)
			continue
		}

		for _, path := range matches {
			if err := a.processLogFile(path); err != nil {
				log.Printf("Error processing log file %s: %v", path, err)
			}
		}
	}
}

// processLogFile reads and processes a single log file
func (a *Agent) processLogFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Skip if file is empty
	if info.Size() == 0 {
		return nil
	}

	// Read last N lines (tail behavior)
	return a.tailLogFile(file, path)
}

// tailLogFile reads the last few lines from a log file
func (a *Agent) tailLogFile(file *os.File, path string) error {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, a.MaxLineLength), a.MaxLineLength)

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		a.processLogLine(path, line)
		lineCount++

		if lineCount >= a.BatchSize {
			break
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// processLogLine processes a single log line
func (a *Agent) processLogLine(source, line string) {
	// Extract metadata and send entry to the aggregator channel
	logEntry := logaggregator.LogEntry{
		Timestamp: time.Now(),
		NodeName:  a.Meta.NodeName,
		Source:    source,
		Message:   line,
		Level:     "info",
	}
	a.logChan <- logEntry

	metrics.LogLines.With(metrics.MakeLabelsForLogLine("", logEntry.NodeName)).Add(1) // TODO: filename (sanitized)
}

// closeLogFiles closes all open log files
func (a *Agent) closeLogFiles() {
	a.logMutex.Lock()
	defer a.logMutex.Unlock()

	for path, file := range a.logFiles {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file %s: %v", path, err)
		}
	}

	a.logFiles = make(map[string]*os.File)
}
