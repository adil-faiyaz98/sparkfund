.PHONY: build run test clean migrate-up migrate-down lint mock dev dev-down fmt tidy check help docker-build docker-push

# Build variables
BINARY_NAME = email-service
MAIN_FILE = cmd/main.go
DOCKER_IMAGE_NAME = your-docker-registry/email-service
DOCKER_IMAGE_TAG = latest

# Go commands
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOCLEAN = $(GOCMD) clean
GOMOD = $(GOCMD) mod
GOTIDY = $(GOCMD) mod tidy
GOVET = $(GOCMD) vet
GOLINT = golangci-lint

# Database variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= email_service

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run database migrations up
migrate-up:
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Run linter
lint:
	$(GOLINT) run

# Generate mocks
mock:
	mockgen -source=internal/repositories/repository.go -destination=internal/repositories/mock_repository.go

# Start development environment
dev:
	docker-compose up -d

# Stop development environment
dev-down:
	docker-compose down

# Format code
fmt:
	$(GOCMD) fmt ./...

# Tidy go modules
tidy:
	$(GOTIDY)

# Run all checks
check: fmt lint test

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

# Push Docker image
docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  migrate-up   - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"
	@echo "  lint         - Run linter"
	@echo "  mock         - Generate mocks"
	@echo "  dev          - Start development environment"
	@echo "  dev-down     - Stop development environment"
	@echo "  fmt          - Format code"
	@echo "  tidy         - Tidy go modules"
	@echo "  check        - Run all checks"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-push  - Push Docker image"