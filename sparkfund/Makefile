.PHONY: build test deploy clean

# Build all services
build:
	@echo "Building services..."
	docker-compose build

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Deploy to Kubernetes
deploy:
	@echo "Deploying to Kubernetes..."
	helm upgrade --install aml-service ./helm/aml-service --namespace sparkfund
	helm upgrade --install kyc-service ./helm/kyc-service --namespace sparkfund

# Clean up
clean:
	@echo "Cleaning up..."
	docker-compose down -v
	rm -rf dist/

# Generate protobuf
proto:
	@echo "Generating protobuf..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# Run locally
run:
	@echo "Running services locally..."
	docker-compose up

# Run with hot reload
dev:
	@echo "Running services with hot reload..."
	docker-compose -f docker-compose.dev.yml up

# Check code style
lint:
	@echo "Running linter..."
	golangci-lint run

# Generate swagger docs
swagger:
	@echo "Generating swagger documentation..."
	swag init -g cmd/main.go

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

# Run security scan
security:
	@echo "Running security scan..."
	gosec ./...
	trivy image sparkfund/aml-service:latest
	trivy image sparkfund/kyc-service:latest