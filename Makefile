.PHONY: build test clean run fmt vet

# Build variables
BINARY_NAME=logging_agent
BUILD_DIR=bin

# Build the application
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/logging_agent

# Run the application
run:
	@echo "Running..."
	@go run ./cmd/logging_agent/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run all checks
check: fmt vet test

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/logging_agent
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/logging_agent
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/logging_agent
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/logging_agent
