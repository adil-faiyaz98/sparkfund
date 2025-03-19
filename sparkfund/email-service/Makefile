.PHONY: build run test clean migrate-up migrate-down lint mock dev dev-down

# Build variables
BINARY_NAME=email-service
MAIN_FILE=cmd/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOCMD) mod tidy
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)
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
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/email_service?sslmode=disable" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/email_service?sslmode=disable" down

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