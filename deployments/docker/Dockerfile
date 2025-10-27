# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o logging_agent \
    ./cmd/logging_agent

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/logging_agent .

# Copy config directory structure (if needed)
RUN mkdir -p /etc/logging-agent

# Note: Running as root is required for accessing /var/log/containers/ on Kubernetes nodes
# The Kubernetes securityContext will handle user permissions

# Expose health check port
EXPOSE 8080

# Set default command
ENTRYPOINT ["/app/logging_agent"]
CMD ["--config", "/etc/logging-agent/config.yaml"]
