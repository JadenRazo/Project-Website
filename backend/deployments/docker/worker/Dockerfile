FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -a -tags netgo,osusergo \
    -o /build/worker ./cmd/worker

# Create minimal production image
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Add non-root user
RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -s /bin/sh -D appuser && \
    mkdir -p /app/config /app/logs && \
    chown -R appuser:appuser /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

# Copy the binary from builder
COPY --from=builder --chown=appuser:appuser /build/worker /app/worker

# Copy configuration files
COPY --chown=appuser:appuser ./config/production.yaml /app/config/production.yaml

# Switch to non-root user
USER appuser

# Set environment variables
ENV ENVIRONMENT=production \
    TZ=UTC \
    CONFIG_PATH=/app/config/production.yaml

# Run the worker
CMD ["/app/worker"]
