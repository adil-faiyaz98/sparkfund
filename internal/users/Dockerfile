# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/users-service ./internal/users/cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/users-service .

# Copy migrations
COPY --from=builder /app/internal/users/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./users-service"] 