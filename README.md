# logging_agent

A simple logging agent written in Go.

## Project Structure

```
.
├── cmd/
│   └── logging_agent/      # Main application entry point
│       └── main.go
├── pkg/
│   └── agent/              # Public packages
│       ├── agent.go
│       └── agent_test.go
├── internal/
│   └── config/             # Private application code
│       └── config.go
├── bin/                    # Compiled binaries (gitignored)
├── go.mod                  # Go module file
├── main.go                 # Alternative simple entry point
├── Makefile                # Build automation
├── .gitignore
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Building

```bash
# Build the application
make build

# Or use go directly
go build -o bin/logging_agent ./cmd/logging_agent
```

### Running

```bash
# Run using make
make run

# Or run directly
go run ./cmd/logging_agent/main.go

# Or run the simple main.go
go run main.go
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Development

```bash
# Format code
make fmt

# Run linter
make vet

# Run all checks
make check

# Clean build artifacts
make clean
```

## Configuration

The application can be configured using environment variables:

- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Logging level (default: info)
- `DEBUG` - Enable debug mode (default: false)

Example:
```bash
PORT=9090 LOG_LEVEL=debug go run ./cmd/logging_agent/main.go
```

## License

MIT