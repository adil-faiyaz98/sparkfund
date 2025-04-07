# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Development stage
FROM golang:1.23-alpine AS development

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

# Install development tools
RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080
EXPOSE 2345

CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.19 AS production

WORKDIR /app

# Install CA certificates and create non-root user
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser

# Copy binary and config from builder
COPY --from=builder /app/bin/kyc-service .
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env.example .env

# Create document storage directory
RUN mkdir -p /data/documents && \
    chown -R appuser:appuser /data/documents && \
    chmod 755 /data/documents

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./kyc-service"]