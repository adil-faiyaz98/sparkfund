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

# Run tests
.PHONY: test
test: unit-test integration-test

# Run unit tests only
.PHONY: unit-test
unit-test:
	$(GOTEST) -short ./...

# Run integration tests only
.PHONY: integration-test
integration-test:
	$(GOTEST) -run Integration ./...

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

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              : Run lint, test, and build"
	@echo "  build            : Build the application"
	@echo "  clean            : Remove build artifacts"
	@echo "  run              : Build and run the application"
	@echo "  test             : Run tests"
	@echo "  test-coverage    : Run tests with coverage"
	@echo "  lint             : Lint the code"
	@echo "  fmt              : Format the code"
	@echo "  docker-build     : Build Docker image"
	@echo "  docker-push      : Push Docker image to registry"
	@echo "  docker-tag-ENV   : Tag and push Docker image for environment (dev, staging, prod)"
	@echo "  deploy-ENV       : Deploy to specific environment (dev, staging, prod)"
	@echo "  view-ENV         : View Kubernetes resources in specific environment"
	@echo "  release-ENV      : Full release pipeline for environment (dev, staging, prod)"
	@echo "  manifests-ENV    : Generate Kubernetes manifests for environment"
	@echo "  help             : Show this help message"
