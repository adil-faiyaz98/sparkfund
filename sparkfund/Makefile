.PHONY: build run test clean proto docker-up docker-down lint

all: build

# Build all services
build:
    @echo "Building all services..."
    @for dir in api-gateway kyc-service aml-service fraud-detection-service credit-scoring-service risk-management-service notification-service consent-management-service logging-service security-service email-service; do \
        echo "Building $$dir..."; \
        (cd $$dir && go build -o bin/service ./cmd/main.go); \
    done

# Run tests for all services
test:
    @echo "Running tests for all services..."
    @for dir in api-gateway kyc-service aml-service fraud-detection-service credit-scoring-service risk-management-service notification-service consent-management-service logging-service security-service email-service; do \
        echo "Testing $$dir..."; \
        (cd $$dir && go test ./...); \
    done

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    @for dir in api-gateway kyc-service aml-service fraud-detection-service credit-scoring-service risk-management-service notification-service consent-management-service logging-service security-service email-service; do \
        echo "Cleaning $$dir..."; \
        (cd $$dir && rm -rf bin); \
    done

# Start all services with Docker Compose
docker-up:
    docker-compose up -d

# Stop all services
docker-down:
    docker-compose down

# Run linting on all services
lint:
    @echo "Running linters..."
    @for dir in api-gateway kyc-service aml-service fraud-detection-service credit-scoring-service risk-management-service notification-service consent-management-service logging-service security-service email-service; do \
        echo "Linting $$dir..."; \
        (cd $$dir && golangci-lint run); \
    done

# Create a new service from template
new-service:
    @read -p "Enter service name: " name; \
    mkdir -p $$name/cmd $$name/internal/config $$name/internal/handlers $$name/internal/services $$name/internal/repositories $$name/internal/models $$name/internal/utils $$name/pkg; \
    cp templates/main.go.tmpl $$name/cmd/main.go; \
    cp templates/go.mod.tmpl $$name/go.mod; \
    cp templates/Dockerfile.tmpl $$name/Dockerfile; \
    sed -i "s/{{SERVICE_NAME}}/$$name/g" $$name/cmd/main.go $$name/go.mod; \
    echo "Created new service: $$name"