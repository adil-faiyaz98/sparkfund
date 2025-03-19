.PHONY: proto build run test clean deps lint

# Proto generation
proto:
	@echo "Generating protobuf code..."
	protoc -I ../../proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		../../proto/accounts/v1/account.proto

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Build the service
build: proto
	@echo "Building service..."
	go build -o bin/accounts-service cmd/main.go

# Run the service
run: build
	@echo "Running service..."
	./bin/accounts-service

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean
	rm -f proto/accounts/v1/*.pb.go

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Generate all (deps, proto, build)
all: deps proto build 