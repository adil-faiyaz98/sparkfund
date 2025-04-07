.PHONY: all build test deploy clean services docker-all docker-push-all proto-all fmt-all lint-all

# List of all services
SERVICES := kyc-service auth-service aml-service investment-service api-gateway security-monitoring-service user-service

# Variables
REGISTRY := sparkfund
VERSION := $(shell git describe --tags --always --dirty)
COMMIT_HASH := $(shell git rev-parse --short HEAD)

# Build all services
build:
    @for service in $(SERVICES); do \
        echo "Building $$service..."; \
        $(MAKE) -C $$service build; \
    done

# Test all services
test:
    @for service in $(SERVICES); do \
        echo "Testing $$service..."; \
        $(MAKE) -C $$service test; \
    done

# Docker operations
docker-all: docker-build docker-push

docker-build:
    @for service in $(SERVICES); do \
        echo "Building Docker image for $$service..."; \
        docker build -t $(REGISTRY)/$$service:$(VERSION) ./$$service; \
        docker tag $(REGISTRY)/$$service:$(VERSION) $(REGISTRY)/$$service:latest; \
    done

docker-push:
    @for service in $(SERVICES); do \
        echo "Pushing Docker image for $$service..."; \
        docker push $(REGISTRY)/$$service:$(VERSION); \
        docker push $(REGISTRY)/$$service:latest; \
    done

# Development tools
fmt-all:
    @for service in $(SERVICES); do \
        echo "Formatting $$service..."; \
        $(MAKE) -C $$service fmt; \
    done

lint-all:
    @for service in $(SERVICES); do \
        echo "Linting $$service..."; \
        $(MAKE) -C $$service lint; \
    done

# Deployment
deploy:
    @for service in $(SERVICES); do \
        echo "Deploying $$service..."; \
        $(MAKE) -C $$service deploy; \
    done

# Clean up
clean:
    @for service in $(SERVICES); do \
        echo "Cleaning $$service..."; \
        $(MAKE) -C $$service clean; \
    done

# List services
services:
    @echo "Available services:"
    @for service in $(SERVICES); do \
        echo "- $$service"; \
    done