# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /build

# Install security scanning tools
RUN apk add --no-cache ca-certificates git && \
    wget -O /usr/local/bin/gosec https://github.com/securego/gosec/releases/download/v2.16.0/gosec_2.16.0_linux_amd64 && \
    chmod +x /usr/local/bin/gosec

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Run security scan
RUN gosec -quiet -exclude-dir=mocks -exclude-dir=test -exclude-dir=docs ./...

# Generate Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.1
RUN swag init -g cmd/main.go --parseDependency --parseInternal

# Run tests
RUN go test -v ./...

# Build with optimizations, trimming, and version info
ARG VERSION=development
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.version=${VERSION}" -o investment-service ./cmd

# Distroless image for ultra-secure minimal runtime
FROM gcr.io/distroless/static-debian11:nonroot AS production

# Copy binary and config from build stage
WORKDIR /app
COPY --from=builder /build/investment-service .
COPY --from=builder /build/docs ./docs
COPY --from=builder /build/config/config.production.yaml ./config/config.yaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set runtime environment
ENV GIN_MODE=release
ENV APP_ENV=production

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/investment-service", "health"]

# Run binary
ENTRYPOINT ["/app/investment-service"]