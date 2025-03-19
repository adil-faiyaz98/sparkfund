#!/bin/bash

# Root directory of the project
ROOT_DIR="."

# List of services
SERVICES=(
    "api-gateway"
    "kyc-service"
    "aml-service"
    "fraud-detection-service"
    "credit-scoring-service"
    "risk-management-service"
    "notification-service"
    "consent-management-service"
    "logging-service"
    "security-service"
    "email-service"
    "blockchain-service"
)

# Create main project directory
mkdir -p "$ROOT_DIR"

# Create each service
for SERVICE in "${SERVICES[@]}"; do
    echo "Creating $SERVICE..."
    
    # Create service directory
    mkdir -p "$ROOT_DIR/$SERVICE"
    
    # Create main.go
    cat > "$ROOT_DIR/$SERVICE/main.go" << EOF
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    
    // Add routes here
    
    log.Printf("Starting $SERVICE on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
EOF

    # Create go.mod
    cat > "$ROOT_DIR/$SERVICE/go.mod" << EOF
module github.com/adil-faiyaz98/sparkfund/$SERVICE

go 1.20

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
)
EOF

    # Create internal directory
    mkdir -p "$ROOT_DIR/$SERVICE/internal"
    
    # Create config directory
    mkdir -p "$ROOT_DIR/$SERVICE/config"
    
    # Create pkg directory
    mkdir -p "$ROOT_DIR/$SERVICE/pkg"
    
    # Create Dockerfile
    cat > "$ROOT_DIR/$SERVICE/Dockerfile" << EOF
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
EOF

    # Create .gitignore
    cat > "$ROOT_DIR/$SERVICE/.gitignore" << EOF
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# IDE specific files
.idea/
.vscode/
*.swp
*.swo

# OS specific files
.DS_Store
Thumbs.db
EOF
done

# Create docker-compose.yml
cat > "$ROOT_DIR/docker-compose.yml" << EOF
version: '3.8'

services:
  api-gateway:
    build: ./api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - kyc-service
      - aml-service
      - fraud-detection-service
      - credit-scoring-service
      - risk-management-service
      - notification-service
      - consent-management-service
      - logging-service
      - security-service
      - email-service
      - blockchain-service

  kyc-service:
    build: ./kyc-service
    ports:
      - "8081:8080"

  aml-service:
    build: ./aml-service
    ports:
      - "8082:8080"

  fraud-detection-service:
    build: ./fraud-detection-service
    ports:
      - "8083:8080"

  credit-scoring-service:
    build: ./credit-scoring-service
    ports:
      - "8084:8080"

  risk-management-service:
    build: ./risk-management-service
    ports:
      - "8085:8080"

  notification-service:
    build: ./notification-service
    ports:
      - "8086:8080"

  consent-management-service:
    build: ./consent-management-service
    ports:
      - "8087:8080"

  logging-service:
    build: ./logging-service
    ports:
      - "8088:8080"

  security-service:
    build: ./security-service
    ports:
      - "8089:8080"

  email-service:
    build: ./email-service
    ports:
      - "8090:8080"

  blockchain-service:
    build: ./blockchain-service
    ports:
      - "8091:8080"
EOF

# Create Makefile
cat > "$ROOT_DIR/Makefile" << EOF
.PHONY: build run test clean

build:
	@for service in \$(shell ls -d */ | grep -v '^\./'); do \
		echo "Building \$${service%/}..."; \
		cd \$${service%/} && go build -o main . && cd ..; \
	done

run:
	docker-compose up

test:
	@for service in \$(shell ls -d */ | grep -v '^\./'); do \
		echo "Testing \$${service%/}..."; \
		cd \$${service%/} && go test ./... && cd ..; \
	done

clean:
	@for service in \$(shell ls -d */ | grep -v '^\./'); do \
		echo "Cleaning \$${service%/}..."; \
		cd \$${service%/} && rm -f main && cd ..; \
	done
EOF

# Create README.md
cat > "$ROOT_DIR/README.md" << EOF
# SparkFund Microservices Platform

This repository contains a collection of microservices for the SparkFund platform.

## Services

- API Gateway
- KYC Service
- AML Service
- Fraud Detection Service
- Credit Scoring Service
- Risk Management Service
- Notification Service
- Consent Management Service
- Logging Service
- Security Service
- Email Service
- Blockchain Service

## Getting Started

### Prerequisites

- Go 1.20 or later
- Docker and Docker Compose
- Make

### Building

\`\`\`bash
make build
\`\`\`

### Running

\`\`\`bash
make run
\`\`\`

### Testing

\`\`\`bash
make test
\`\`\`

### Cleaning

\`\`\`bash
make clean
\`\`\`

## Architecture

Each service is a standalone Go application that communicates with other services through HTTP APIs. The API Gateway serves as the entry point for all external requests.

## License

This project is licensed under the MIT License.
EOF

echo "Repository structure created successfully!" 