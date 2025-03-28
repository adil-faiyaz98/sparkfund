.PHONY: all build test deploy clean services

# List of all services
SERVICES := kyc-service auth-service aml-service credit-scoring-service investment-service

# Build all services
all: build

# Build all services
build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && make build && cd ..; \
	done

# Test all services
test:
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $$service && make test && cd ..; \
	done

# Deploy all services
deploy:
	@for service in $(SERVICES); do \
		echo "Deploying $$service..."; \
		cd $$service && make deploy && cd ..; \
	done

# Clean all services
clean:
	@for service in $(SERVICES); do \
		echo "Cleaning $$service..."; \
		cd $$service && make clean && cd ..; \
	done

# List all services
services:
	@echo "Available services:"
	@for service in $(SERVICES); do \
		echo "- $$service"; \
	done

# Build and push Docker images for all services
docker-all:
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		cd $$service && make docker-build && cd ..; \
	done

# Push Docker images for all services
docker-push-all:
	@for service in $(SERVICES); do \
		echo "Pushing Docker image for $$service..."; \
		cd $$service && make docker-push && cd ..; \
	done

# Generate protobuf for all services
proto-all:
	@for service in $(SERVICES); do \
		echo "Generating protobuf for $$service..."; \
		cd $$service && make proto && cd ..; \
	done

# Format code for all services
fmt-all:
	@for service in $(SERVICES); do \
		echo "Formatting code for $$service..."; \
		cd $$service && make fmt && cd ..; \
	done

# Lint code for all services
lint-all:
	@for service in $(SERVICES); do \
		echo "Linting code for $$service..."; \
		cd $$service && make lint && cd ..; \
	done 