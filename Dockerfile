# Multi-stage build for Go application
# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o logging-agent \
    ./cmd/logging_agent

# Stage 2: Create minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/logging-agent /usr/local/bin/logging-agent

# Create non-root user structure (for metadata only, as scratch has no user management)
# The securityContext in Kubernetes will handle the actual user

# Set working directory
WORKDIR /

# Expose health check port
EXPOSE 8080

# Run as non-root user (specified in K8s manifest)
USER 1000

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/logging-agent"]
