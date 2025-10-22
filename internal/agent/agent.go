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
	"github.com/Semerokozlyat/logging_agent/internal/pkg/metrics"
)

// Agent represents the logging agent
type Agent struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	logFiles map[string]*os.File
	logMutex sync.RWMutex

	healthCheckServer *http.Server

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

// New creates a new Agent instance
func New(cfg *config.Config) *Agent {
	ctx, cancel := context.WithCancel(context.Background())

	// Get configuration from environment
	nodeName := getEnv("NODE_NAME", "unknown-node")
	podName := getEnv("POD_NAME", "logging-agent")
	namespace := getEnv("POD_NAMESPACE", "logging-system")
	logLevel := getEnv("LOG_LEVEL", "info")

	// Init metrics
	metrics.InitMetricsCollector()

	log.Printf("Initializing logging agent with configuration: %+v", cfg)

	return &Agent{
		LogLevel:   logLevel,
		OutputPath: "/var/log/agent.log",
		LogPaths: []string{
			"/var/log/containers/*.log",
			"/var/log/pods/**/*.log",
		},
		CollectionInterval: 10 * time.Second,
		BatchSize:          100,
		MaxLineLength:      16384,
		NodeName:           nodeName,
		PodName:            podName,
		Namespace:          namespace,
		ctx:                ctx,
		cancel:             cancel,
		logFiles:           make(map[string]*os.File),

		healthCheckServer: httpserver.NewHealthCheckServer(&cfg.HTTPServer),
	}
}

// Run starts the agent
func (a *Agent) Run() error {
	log.Printf("Starting Logging Agent on node: %s", a.NodeName)

	// Start health check server
	a.wg.Add(1)
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

	// Start log collection
	a.wg.Add(1)
	go a.collectLogs()

	// Wait for context cancellation
	<-a.ctx.Done()

	log.Println("Shutting down agent...")
	a.wg.Wait()

	// Close all open log files
	a.closeLogFiles()

	log.Println("Agent stopped gracefully")
	return nil
}

// Stop gracefully stops the agent
func (a *Agent) Stop() {
	log.Println("Stop signal received")
	a.cancel()
}

// collectLogs collects logs from specified paths
func (a *Agent) collectLogs() {
	defer a.wg.Done()

	ticker := time.NewTicker(a.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
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
	// Extract metadata from log line
	logEntry := LogEntry{
		Timestamp: time.Now(),
		NodeName:  a.NodeName,
		Source:    source,
		Message:   line,
		Level:     "info",
	}

	// Output log entry (to stdout for now)
	a.outputLog(logEntry)
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	NodeName  string    `json:"node_name"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

// outputLog outputs a log entry
func (a *Agent) outputLog(entry LogEntry) {
	// Format as JSON-like output
	fmt.Printf("[%s] [%s] [%s] %s: %s\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.NodeName,
		entry.Source,
		entry.Message,
	)
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

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
