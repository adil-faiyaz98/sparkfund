# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always || echo 'dev') -X main.buildTime=$(date +%Y-%m-%dT%H:%M:%S)" \
    -a -installsuffix cgo -o service ./cmd/main.go

# Development stage
FROM golang:1.21-alpine AS development

# Install development tools
RUN apk add --no-cache git ca-certificates tzdata curl && \
    update-ca-certificates && \
    go install github.com/cosmtrek/air@latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose ports
EXPOSE 8080 9090

# Set entrypoint
ENTRYPOINT ["air", "-c", ".air.toml"]

# Production stage
FROM gcr.io/distroless/static:nonroot AS production

# Copy binary and config
WORKDIR /
COPY --from=builder /app/service /service
COPY --from=builder /app/config /config
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use non-root user
USER nonroot:nonroot

# Expose ports
EXPOSE 8080 9090

# Set health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD [ "/service", "health", "check" ] || exit 1

# Run the binary
ENTRYPOINT ["/service"]
