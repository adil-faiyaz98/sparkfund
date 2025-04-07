# Service Makefile

# Variables
SERVICE_NAME := service-template
BINARY_NAME := $(SERVICE_NAME)
DOCKER_IMAGE := sparkfund/$(SERVICE_NAME)
DOCKER_TAG := latest
MAIN_PATH := ./cmd/main.go
MIGRATION_PATH := ./internal/database/migrations

# Go commands
GO := go
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

# Build settings
BUILD_DIR := build
LDFLAGS := -ldflags "-w -s -X main.version=$(shell git describe --tags --always) -X main.buildTime=$(shell date -u '+%Y-%m-%d %H:%M:%S')"

# Docker commands
DOCKER := docker
DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_COMPOSE_DEV_FILE := docker-compose.dev.yml

# Kubernetes
KUBECTL := kubectl
K8S_NAMESPACE := sparkfund

# Database
DB_DSN := postgres://postgres:postgres@localhost:5432/$(SERVICE_NAME)?sslmode=disable

# Build the service
.PHONY: build
build:
	CGO_ENABLED=0 $(GO) build -a -installsuffix cgo $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the service
.PHONY: run
run:
	$(GO) run $(MAIN_PATH)

# Run with hot reload
.PHONY: dev
dev:
	air -c .air.toml

# Build Docker image
.PHONY: docker-build
docker-build:
	$(DOCKER) build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Push Docker image
.PHONY: docker-push
docker-push:
	$(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run Docker container
.PHONY: docker-run
docker-run:
	$(DOCKER) run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Start Docker Compose development environment
.PHONY: compose-up
compose-up:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV_FILE) up -d

# Stop Docker Compose development environment
.PHONY: compose-down
compose-down:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV_FILE) down

# Show Docker Compose logs
.PHONY: compose-logs
compose-logs:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV_FILE) logs -f

# Run tests
.PHONY: test
test:
	$(GO) test -v -race -cover ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run linter
.PHONY: lint
lint:
	golangci-lint run ./...

# Run security scanner
.PHONY: security-scan
security-scan:
	gosec -exclude-dir=vendor ./...

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init -g $(MAIN_PATH) -o docs/swagger

# Generate mocks
.PHONY: mocks
mocks:
	go generate ./...

# Format code
.PHONY: fmt
fmt:
	goimports -w .

# Deploy to Kubernetes
.PHONY: k8s-deploy
k8s-deploy:
	$(KUBECTL) apply -f k8s/

# Delete from Kubernetes
.PHONY: k8s-delete
k8s-delete:
	$(KUBECTL) delete -f k8s/

# Run database migrations up
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DB_DSN)" up

# Run database migrations down
.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DB_DSN)" down

# Help
.PHONY: help
help:
	@echo "$(SERVICE_NAME) Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build             Build the service"
	@echo "  clean             Clean build artifacts"
	@echo "  run               Run the service"
	@echo "  dev               Run with hot reload"
	@echo "  docker-build      Build Docker image"
	@echo "  docker-push       Push Docker image"
	@echo "  docker-run        Run Docker container"
	@echo "  compose-up        Start Docker Compose environment"
	@echo "  compose-down      Stop Docker Compose environment"
	@echo "  compose-logs      Show Docker Compose logs"
	@echo "  test              Run tests"
	@echo "  test-coverage     Run tests with coverage"
	@echo "  lint              Run linter"
	@echo "  security-scan     Run security scanner"
	@echo "  swagger           Generate Swagger documentation"
	@echo "  mocks             Generate mocks"
	@echo "  fmt               Format code"
	@echo "  k8s-deploy        Deploy to Kubernetes"
	@echo "  k8s-delete        Delete from Kubernetes"
	@echo "  migrate-up        Run database migrations"
	@echo "  migrate-down      Rollback database migrations"
	@echo "  help              Show this help"
