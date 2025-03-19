.PHONY: build run test clean migrate-up migrate-down lint mock dev dev-down fmt check help

# Build variables
BINARY_NAME=email-service
MAIN_FILE=cmd/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOTIDY=$(GOCMD) mod tidy

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

# Run database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/email_service?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/email_service?sslmode=disable" down

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

# Run all checks
check: fmt lint test

# Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  migrate-up   - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"
	@echo "  mock         - Generate mocks"
	@echo "  dev          - Start development environment"
	@echo "  dev-down     - Stop development environment"
	@echo "  fmt          - Format code"
	@echo "  check        - Run all checks (fmt, lint, test)"
	@echo "  help         - Show this help message"