# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o bin/user-service ./cmd/user-service

# Production stage
FROM alpine:3.19 AS production

WORKDIR /app

# Install CA certificates and create non-root user
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser

# Copy binary and config from builder
COPY --from=builder /app/bin/user-service .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8084

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8084/health || exit 1

CMD ["./user-service"]
