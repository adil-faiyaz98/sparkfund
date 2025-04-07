# Variables
SERVICE_NAME := kyc-service
REGISTRY := sparkfund
VERSION := $(shell git describe --tags --always --dirty)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go related variables
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test

# Build flags
LDFLAGS := -w -s \
    -X main.version=$(VERSION) \
    -X main.commitHash=$(COMMIT_HASH) \
    -X main.buildTime=$(BUILD_TIME)

.PHONY: all build test clean fmt lint docker-build docker-push run mock help generate docs deps

all: clean deps generate test build

build:
    @echo "Building $(SERVICE_NAME)..."
    $(GOBUILD) -ldflags "$(LDFLAGS)" -o bin/$(SERVICE_NAME) ./cmd/api

# Test targets
.PHONY: test test-unit test-integration test-e2e

test-unit:
    @echo "Running unit tests..."
    go test -v ./tests/unit/...

test-integration:
    @echo "Running integration tests..."
    go test -v ./tests/integration/...

test-e2e:
    @echo "Running e2e tests..."
    go test -v ./tests/e2e/...

test: test-unit test-integration test-e2e

lint:
    @echo "Running linter..."
    golangci-lint run --timeout=5m

generate:
    @echo "Generating code..."
    $(GOCMD) generate ./...
    swag init -g cmd/api/main.go --parseDependency --parseInternal

docker-build:
    docker build -t $(REGISTRY)/$(SERVICE_NAME):$(VERSION) \
        --build-arg VERSION=$(VERSION) \
        --build-arg COMMIT_HASH=$(COMMIT_HASH) \
        --build-arg BUILD_TIME=$(BUILD_TIME) \
        -f deployments/docker/Dockerfile .

dev:
    @echo "Starting development environment..."
    @chmod +x scripts/setup-dev.sh
    @./scripts/setup-dev.sh

# Development commands
.PHONY: dev dev-down dev-logs dev-clean

# Start development environment
dev:
    @echo "Starting development environment..."
    @chmod +x scripts/setup-dev.sh
    @./scripts/setup-dev.sh

# Stop development environment
dev-down:
    @echo "Stopping development environment..."
    docker-compose -f docker-compose.dev.yml down

# Show development logs
dev-logs:
    @echo "Showing development logs..."
    docker-compose -f docker-compose.dev.yml logs -f

# Clean development environment
dev-clean:
    @echo "Cleaning development environment..."
    docker-compose -f docker-compose.dev.yml down -v
    rm -rf tmp/
    rm -rf data/documents/*

# Generate mocks
mock:
    @echo "Generating mocks..."
    mockgen -source=internal/repository/kyc.go -destination=internal/repository/mock/kyc_mock.go
    mockgen -source=internal/service/document_service.go -destination=internal/service/mock/document_service_mock.go

# Generate API documentation
docs:
    @echo "Generating API documentation..."
    swag init -g cmd/api/main.go --parseDependency --parseInternal

# Download dependencies
deps:
    @echo "Downloading dependencies..."
    $(GOMOD) download
    $(GOMOD) tidy

# Generate code (protobuf, mocks, etc.)
generate:
    @echo "Generating code..."
    $(GOCMD) generate ./...

# Docker compose commands
compose-up:
    @echo "Starting docker compose..."
    $(DOCKER_COMPOSE) up -d

compose-down:
    @echo "Stopping docker compose..."
    $(DOCKER_COMPOSE) down

compose-logs:
    @echo "Showing logs..."
    $(DOCKER_COMPOSE) logs -f

# Database commands
db-migrate:
    @echo "Running database migrations..."
    go run cmd/api/main.go migrate

db-rollback:
    @echo "Rolling back database migrations..."
    go run cmd/api/main.go migrate rollback

# Security scan
security-scan:
    @echo "Running security scan..."
    gosec -quiet ./...
    nancy go.sum

# Show help
help:
    @echo "Available targets:"
    @echo "  build         - Build the service"
    @echo "  test          - Run tests with coverage"
    @echo "  clean         - Clean build artifacts"
    @echo "  fmt           - Format code"
    @echo "  lint          - Run linter"
    @echo "  docker-build  - Build Docker image"
    @echo "  docker-push   - Push Docker image"
    @echo "  run           - Run service locally"
    @echo "  mock          - Generate mocks"
    @echo "  docs          - Generate API documentation"
    @echo "  deps          - Download dependencies"
    @echo "  generate      - Generate code"
    @echo "  compose-up    - Start docker compose"
    @echo "  compose-down  - Stop docker compose"
    @echo "  compose-logs  - Show docker compose logs"
    @echo "  db-migrate    - Run database migrations"
    @echo "  db-rollback   - Rollback database migrations"
    @echo "  security-scan - Run security scan"
