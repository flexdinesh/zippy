# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o zippy \
    ./cmd/zippy

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 zippy && \
    adduser -D -u 1000 -G zippy zippy

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/zippy /app/zippy

# Create directories for backups and data
RUN mkdir -p /data /backups && \
    chown -R zippy:zippy /app /data /backups

# Switch to non-root user
USER zippy

# Default configuration path
ENV CONFIG_PATH=/app/config.yaml

ENTRYPOINT ["/app/zippy"]
CMD ["-config", "/app/config.yaml"]
