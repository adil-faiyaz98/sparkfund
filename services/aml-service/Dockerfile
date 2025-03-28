# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies and security updates
RUN apk add --no-cache git && \
    apk upgrade --no-cache

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o aml-service ./services/aml-service/cmd/main.go

# Security scan stage
FROM aquasec/trivy:latest AS security-scan
COPY --from=builder /app/aml-service /aml-service
RUN trivy filesystem --no-progress --ignore-unfixed --severity HIGH,CRITICAL /aml-service

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install security updates and create non-root user
RUN apk upgrade --no-cache && \
    adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

# Copy binary from builder
COPY --from=builder /app/aml-service .

# Copy config files
COPY --from=builder /app/config ./config

# Set security context
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./aml-service"] 