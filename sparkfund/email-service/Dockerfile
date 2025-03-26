# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/email-service ./cmd/main.go

# Final stage
FROM alpine:3.17

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone
ENV TZ=Etc/UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create user and group
RUN addgroup -S app && adduser -S -G app app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/email-service .

# Copy migrations
COPY migrations ./migrations

# Change ownership
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD wget --quiet --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./email-service"]