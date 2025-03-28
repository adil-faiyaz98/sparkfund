.PHONY: build test deploy clean

# Variables
SERVICE_NAME := auth-service
DOCKER_IMAGE := sparkfund/$(SERVICE_NAME)
VERSION := $(shell git describe --tags --always --dirty)
K8S_NAMESPACE := sparkfund

# Build the service
build:
	go build -o $(SERVICE_NAME) ./cmd/main.go

# Run tests
test:
	go test -v ./...

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

# Push Docker image
docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

# Deploy to Kubernetes
deploy:
	kubectl apply -k k8s/

# Clean up
clean:
	rm -f $(SERVICE_NAME)
	docker rmi $(DOCKER_IMAGE):$(VERSION) || true
	docker rmi $(DOCKER_IMAGE):latest || true

# Generate protobuf
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto

# Run locally
run:
	go run ./cmd/main.go

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate mocks
mock:
	mockgen -source=internal/service.go -destination=internal/mocks/service_mock.go 