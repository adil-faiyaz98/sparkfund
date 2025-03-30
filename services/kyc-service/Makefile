.PHONY: build test clean fmt lint docker-build docker-push run mock

# Service name
SERVICE_NAME := kyc-service

# Docker registry
REGISTRY := sparkfund

# Build the service
build:
	go build -o bin/$(SERVICE_NAME) ./cmd/api

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t $(REGISTRY)/$(SERVICE_NAME):latest .

# Push Docker image
docker-push:
	docker push $(REGISTRY)/$(SERVICE_NAME):latest

# Run the service locally
run:
	go run ./cmd/api

# Generate mocks
mock:
	mockgen -source=internal/repository/kyc.go -destination=internal/repository/mock_kyc.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run security scan
security:
	gosec ./...

# Run all checks
check: fmt lint test security

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the service"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-push  - Push Docker image"
	@echo "  run          - Run the service locally"
	@echo "  mock         - Generate mocks"
	@echo "  deps         - Install dependencies"
	@echo "  security     - Run security scan"
	@echo "  check        - Run all checks"
	@echo "  help         - Show this help message" 