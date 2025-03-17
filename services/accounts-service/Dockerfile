FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o accounts-service ./cmd/server

# Use a minimal alpine image for the final container
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/accounts-service .

# Expose ports for HTTP and gRPC
EXPOSE 8080 50051

# Command to run the executable
CMD ["./accounts-service"]