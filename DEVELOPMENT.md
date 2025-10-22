# Development Guide

This guide covers local development workflows for the logging agent.

## Development Environment Setup

### Requirements

- Go 1.21 or higher
- Docker
- kubectl (for testing Kubernetes deployments)
- Make
- Git

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/Semerokozlyat/logging_agent.git
cd logging_agent

# Install dependencies
go mod download

# Verify everything builds
make build

# Run tests
make test
```

## Development Workflow

### 1. Code Changes

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
vim pkg/agent/agent.go

# Format code
make fmt

# Run linter
make vet

# Run tests
make test
```

### 2. Local Testing

#### Option A: Run Directly

```bash
# Set required environment variables
export NODE_NAME=dev-node
export POD_NAME=logging-agent-dev
export POD_NAMESPACE=default
export LOG_LEVEL=debug

# Run the agent
go run ./cmd/logging_agent/main.go
```

#### Option B: Build and Run Binary

```bash
# Build
make build

# Run
./bin/logging-agent
```

#### Option C: Docker

```bash
# Build development image
docker build -t logging-agent:dev -f Dockerfile.dev .

# Run with mounted volumes
docker run --rm \
  -v /var/log:/var/log:ro \
  -e NODE_NAME=docker-node \
  -e LOG_LEVEL=debug \
  -p 8080:8080 \
  logging-agent:dev
```

### 3. Testing Health Endpoints

```bash
# In another terminal
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

## Project Structure

### Key Directories

```
├── cmd/                    # Application entry points
│   └── logging_agent/      # Main agent binary
├── pkg/                    # Public libraries
│   └── agent/              # Core agent logic
├── internal/               # Private libraries
│   └── config/             # Configuration
├── deployments/            # Kubernetes manifests & scripts
└── docs/                   # Documentation (if added)
```

### Adding New Features

#### 1. Add Output Plugin

Create a new output plugin in `pkg/agent/output/`:

```go
// pkg/agent/output/elasticsearch.go
package output

type ElasticsearchOutput struct {
    url string
}

func NewElasticsearchOutput(url string) *ElasticsearchOutput {
    return &ElasticsearchOutput{url: url}
}

func (e *ElasticsearchOutput) Write(entry LogEntry) error {
    // Implementation
}
```

#### 2. Add Tests

```go
// pkg/agent/output/elasticsearch_test.go
package output

import "testing"

func TestElasticsearchOutput(t *testing.T) {
    // Test implementation
}
```

#### 3. Wire it up

Update `pkg/agent/agent.go` to use the new output.

## Testing

### Unit Tests

```bash
# Run all tests
make test

# Run specific package
go test -v ./pkg/agent

# Run specific test
go test -v ./pkg/agent -run TestAgentCreation

# With coverage
make test-coverage
open coverage.html
```

### Integration Tests

```bash
# Create integration test
mkdir -p test/integration
vim test/integration/agent_test.go

# Run integration tests
go test -v ./test/integration -tags=integration
```

### Kubernetes Tests

```bash
# Deploy to local cluster (Minikube/Kind)
cd deployments
./deploy.sh

# Verify
kubectl get pods -n logging-system

# Check logs
kubectl logs -n logging-system -l app.kubernetes.io/name=logging-agent

# Test health
kubectl port-forward -n logging-system daemonset/logging-agent 8080:8080
curl http://localhost:8080/health
```

## Debugging

### Enable Debug Logging

```bash
export LOG_LEVEL=debug
go run ./cmd/logging_agent/main.go
```

### Use Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/logging_agent/main.go

# Set breakpoint
(dlv) break agent.go:75
(dlv) continue
```

### Profile Performance

```bash
# CPU profiling
go run ./cmd/logging_agent/main.go -cpuprofile=cpu.prof

# Analyze
go tool pprof cpu.prof
```

### Debug in Kubernetes

```bash
# Use development image with debugging tools
docker build -t logging-agent:debug -f Dockerfile.dev .

# Deploy with debug image
# Update daemonset.yaml image to logging-agent:debug

# Exec into pod
kubectl exec -it -n logging-system logging-agent-xxxxx -- sh

# Check logs
kubectl logs -n logging-system logging-agent-xxxxx -f
```

## Code Quality

### Linting

```bash
# Go vet
make vet

# golangci-lint (if installed)
golangci-lint run
```

### Formatting

```bash
# Format all code
make fmt

# Check formatting
gofmt -l .
```

### Security Scanning

```bash
# Scan for vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Container scanning
docker scan logging-agent:latest
```

## Benchmarking

```bash
# Add benchmark tests
vim pkg/agent/agent_bench_test.go

# Run benchmarks
go test -bench=. -benchmem ./pkg/agent

# Profile benchmarks
go test -bench=. -benchmem -cpuprofile=cpu.prof ./pkg/agent
go tool pprof cpu.prof
```

## Release Process

### 1. Version Bump

```bash
# Update version in relevant files
# - README.md
# - Makefile
# - deployments/kubernetes/daemonset.yaml

# Tag release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 2. Build Release Artifacts

```bash
# Build for multiple platforms
make build-all

# Verify binaries
ls -lh bin/
```

### 3. Build and Push Docker Image

```bash
# Build
docker build -t your-registry/logging-agent:v1.0.0 .

# Tag as latest
docker tag your-registry/logging-agent:v1.0.0 your-registry/logging-agent:latest

# Push
docker push your-registry/logging-agent:v1.0.0
docker push your-registry/logging-agent:latest
```

### 4. Update Documentation

```bash
# Update CHANGELOG.md
# Update README.md with new features
# Update deployment guides if needed
```

## Continuous Integration

### GitHub Actions Example

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make test
      - run: make build
```

## Tips & Tricks

### Fast Iteration

```bash
# Use air for hot reload (install first)
go install github.com/cosmtrek/air@latest

# Create .air.toml config
air init

# Run with hot reload
air
```

### Quick Container Test

```bash
# Build and run in one command
docker build -t logging-agent:test . && \
docker run --rm -e NODE_NAME=test logging-agent:test
```

### Test RBAC Changes

```bash
# Apply only RBAC resources
kubectl apply -f deployments/kubernetes/clusterrole.yaml
kubectl apply -f deployments/kubernetes/clusterrolebinding.yaml

# Test permissions
kubectl auth can-i list pods --as=system:serviceaccount:logging-system:logging-agent
```

## Common Issues

### Import Cycle

```
package github.com/Semerokozlyat/logging_agent/pkg/agent
    imports github.com/Semerokozlyat/logging_agent/internal/config
    imports github.com/Semerokozlyat/logging_agent/pkg/agent: import cycle
```

**Solution**: Restructure packages to avoid circular dependencies.

### File Permissions in Docker

```
Error: permission denied reading /var/log/containers
```

**Solution**: Ensure proper volume mounts and security context in DaemonSet.

### Context Deadline Exceeded

```
Error: context deadline exceeded
```

**Solution**: Increase timeouts in health checks or fix blocking operations.
