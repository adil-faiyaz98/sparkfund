# SparkFund Makefile

# Variables
SHELL := /bin/bash
GO := go
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin
GOLANGCI_LINT_VERSION := v1.55.2

# Services
SERVICES := api-gateway kyc-service investment-service user-service

# Docker
DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_FILE := docker-compose.yml

# Kubernetes
KUBECTL := kubectl
K8S_NAMESPACE := sparkfund

# Tools
.PHONY: tools
tools:
	@echo "Installing tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/golang/mock/mockgen@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install github.com/go-delve/delve/cmd/dlv@latest
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO) install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Development
.PHONY: dev
dev:
	$(DOCKER_COMPOSE) up -d

.PHONY: dev-down
dev-down:
	$(DOCKER_COMPOSE) down

.PHONY: dev-logs
dev-logs:
	$(DOCKER_COMPOSE) logs -f

# Production
.PHONY: prod
prod:
	$(DOCKER_COMPOSE) up -d

.PHONY: prod-down
prod-down:
	$(DOCKER_COMPOSE) down

.PHONY: prod-logs
prod-logs:
	$(DOCKER_COMPOSE) logs -f

# Testing
.PHONY: test
test:
	$(GO) test -v -race -cover ./...

.PHONY: test-coverage
test-coverage:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Linting
.PHONY: lint
lint:
	golangci-lint run ./...

# Security scanning
.PHONY: security-scan
security-scan:
	gosec -exclude-dir=vendor ./...

# Build all services
.PHONY: build-all
build-all:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(MAKE) -C services/$$service build; \
	done

# Clean all services
.PHONY: clean-all
clean-all:
	@for service in $(SERVICES); do \
		echo "Cleaning $$service..."; \
		$(MAKE) -C services/$$service clean; \
	done

# Generate API documentation
.PHONY: swagger
swagger:
	@for service in $(SERVICES); do \
		echo "Generating Swagger docs for $$service..."; \
		cd services/$$service && swag init -g cmd/main.go -o docs/swagger; \
	done

# Generate mocks
.PHONY: mocks
mocks:
	@for service in $(SERVICES); do \
		echo "Generating mocks for $$service..."; \
		$(MAKE) -C services/$$service mocks; \
	done

# Format code
.PHONY: fmt
fmt:
	goimports -w .

# Kubernetes deployment
.PHONY: k8s-deploy
k8s-deploy:
	@for service in $(SERVICES); do \
		echo "Deploying $$service to Kubernetes..."; \
		$(KUBECTL) apply -f services/$$service/k8s/; \
	done

.PHONY: k8s-delete
k8s-delete:
	@for service in $(SERVICES); do \
		echo "Deleting $$service from Kubernetes..."; \
		$(KUBECTL) delete -f services/$$service/k8s/; \
	done

# Database migrations
.PHONY: migrate-up
migrate-up:
	@for service in $(SERVICES); do \
		echo "Running migrations for $$service..."; \
		$(MAKE) -C services/$$service migrate-up; \
	done

.PHONY: migrate-down
migrate-down:
	@for service in $(SERVICES); do \
		echo "Rolling back migrations for $$service..."; \
		$(MAKE) -C services/$$service migrate-down; \
	done

# Help
.PHONY: help
help:
	@echo "SparkFund Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  tools             Install development tools"
	@echo "  dev               Start development environment"
	@echo "  dev-down          Stop development environment"
	@echo "  dev-logs          Show development logs"
	@echo "  prod              Start production environment"
	@echo "  prod-down         Stop production environment"
	@echo "  prod-logs         Show production logs"
	@echo "  test              Run tests"
	@echo "  test-coverage     Run tests with coverage"
	@echo "  lint              Run linter"
	@echo "  security-scan     Run security scanner"
	@echo "  build-all         Build all services"
	@echo "  clean-all         Clean all services"
	@echo "  swagger           Generate Swagger documentation"
	@echo "  mocks             Generate mocks"
	@echo "  fmt               Format code"
	@echo "  k8s-deploy        Deploy to Kubernetes"
	@echo "  k8s-delete        Delete from Kubernetes"
	@echo "  migrate-up        Run database migrations"
	@echo "  migrate-down      Rollback database migrations"
	@echo "  help              Show this help"
