# Money Pulse Makefile
# Common commands for development, building, and deployment

# Variables
APP_NAME := money-pulse
DOCKER_REGISTRY := registry.example.com
VERSION := $(shell git describe --tags --always --dirty)
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(APP_NAME)

# Go related variables
GO := go
GOFLAGS := -v
GOTEST := $(GO) test
GOBUILD := $(GO) build $(GOFLAGS)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

# Kubernetes related variables
KUSTOMIZE := kubectl kustomize
KUBECTL := kubectl
ENVIRONMENTS := dev staging prod

# Default target
.PHONY: all
all: lint test build

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

# Run the application
.PHONY: run
run: build
	./bin/$(APP_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	rm -f *.out
	find . -name "*.pb.go" -type f -delete

# Run all tests
.PHONY: test
test: unit-test integration-test e2e-test

# Run unit tests only
.PHONY: unit-test
unit-test:
	@echo "Running unit tests..."
	$(GOTEST) -short -v ./...

# Run integration tests only
.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

# Run end-to-end tests only
.PHONY: e2e-test
e2e-test:
	@echo "Running end-to-end tests..."
	$(GOTEST) -v -tags=e2e ./tests/e2e/...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out

# Run tests continuously on file changes
.PHONY: test-watch
test-watch:
	which goconvey > /dev/null || go install github.com/smartystreets/goconvey@latest
	goconvey

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	gofmt -s -w $(GOFILES)

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .

# Push Docker image to registry
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

# Tag Docker image for specific environment
.PHONY: docker-tag-%
docker-tag-%: docker-build
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):$*
	docker push $(DOCKER_IMAGE):$*

# Deploy to specific environment
.PHONY: deploy-%
deploy-%:
	@if [[ ! " $(ENVIRONMENTS) " =~ " $* " ]]; then \
		echo "Environment '$*' is not valid. Valid environments are: $(ENVIRONMENTS)"; \
		exit 1; \
	fi
	$(KUSTOMIZE) k8s/overlays/$* | $(KUBECTL) apply -f -

# View Kubernetes resources in specific environment
.PHONY: view-%
view-%:
	@if [[ ! " $(ENVIRONMENTS) " =~ " $* " ]]; then \
		echo "Environment '$*' is not valid. Valid environments are: $(ENVIRONMENTS)"; \
		exit 1; \
	fi
	$(KUBECTL) get all -n money-pulse-$*

# Full build and deployment pipeline for specific environment
.PHONY: release-%
release-%: test docker-tag-$* deploy-$*
	@echo "Released to $* environment"

# Generate k8s manifests without applying
.PHONY: manifests-%
manifests-%:
	$(KUSTOMIZE) k8s/overlays/$* > k8s-$*.yaml

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/smartystreets/goconvey@latest

# Generate protobuf code for all services
.PHONY: proto-all
proto-all: proto-accounts proto-users proto-transactions proto-loans proto-investments proto-reports

# Generate protobuf code for accounts service
.PHONY: proto-accounts
proto-accounts:
	@echo "Generating protobuf code for accounts service..."
	protoc -I ./proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./proto/common/v1/common.proto \
		./proto/accounts/v1/account.proto

# Generate protobuf code for users service
.PHONY: proto-users
proto-users:
	@echo "Generating protobuf code for users service..."
	protoc -I ./proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./proto/common/v1/common.proto \
		./proto/users/v1/user.proto

# Generate protobuf code for transactions service
.PHONY: proto-transactions
proto-transactions:
	@echo "Generating protobuf code for transactions service..."
	protoc -I ./proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./proto/common/v1/common.proto \
		./proto/transactions/v1/transaction.proto

# Generate protobuf code for loans service
.PHONY: proto-loans
proto-loans:
	@echo "Generating protobuf code for loans service..."
	protoc -I ./proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./proto/common/v1/common.proto \
		./proto/loans/v1/loan.proto

# Generate protobuf code for investments service
.PHONY: proto-investments
proto-investments:
	@echo "Generating protobuf code for investments service..."
	protoc -I ./proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./proto/common/v1/common.proto \
		./proto/investments/v1/investment.proto

# Generate protobuf code for reports service
.PHONY: proto-reports
proto-reports:
	@echo "Generating protobuf code for reports service..."
	protoc -I ./proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./proto/common/v1/common.proto \
		./proto/reports/v1/report.proto

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              : Run lint, test, and build"
	@echo "  build            : Build the application"
	@echo "  clean            : Remove build artifacts"
	@echo "  run              : Build and run the application"
	@echo "  test             : Run all tests (unit, integration, e2e)"
	@echo "  unit-test        : Run unit tests only"
	@echo "  integration-test : Run integration tests only"
	@echo "  e2e-test         : Run end-to-end tests only"
	@echo "  test-coverage    : Run tests with coverage"
	@echo "  test-watch       : Run tests continuously on file changes"
	@echo "  lint             : Lint the code"
	@echo "  fmt              : Format the code"
	@echo "  docker-build     : Build Docker image"
	@echo "  docker-push      : Push Docker image to registry"
	@echo "  docker-tag-ENV   : Tag and push Docker image for environment (dev, staging, prod)"
	@echo "  deploy-ENV       : Deploy to specific environment (dev, staging, prod)"
	@echo "  view-ENV         : View Kubernetes resources in specific environment"
	@echo "  release-ENV      : Full release pipeline for environment (dev, staging, prod)"
	@echo "  manifests-ENV    : Generate Kubernetes manifests for environment"
	@echo "  deps             : Install dependencies"
	@echo "  proto-all        : Generate protobuf code for all services"
	@echo "  proto-accounts   : Generate protobuf code for accounts service"
	@echo "  proto-users      : Generate protobuf code for users service"
	@echo "  proto-transactions: Generate protobuf code for transactions service"
	@echo "  proto-loans      : Generate protobuf code for loans service"
	@echo "  proto-investments: Generate protobuf code for investments service"
	@echo "  proto-reports    : Generate protobuf code for reports service"
	@echo "  help             : Show this help message"
