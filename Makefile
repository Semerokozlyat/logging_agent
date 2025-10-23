.PHONY: build build-all test test-coverage fmt vet

# Build variables
BINARY_NAME=logging_agent
BUILD_DIR=bin

# Build the application
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/logging_agent

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

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/logging_agent
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/logging_agent
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/logging_agent
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/logging_agent
