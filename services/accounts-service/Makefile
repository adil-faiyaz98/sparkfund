.PHONY: proto build run test clean

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/accounts/v1/account.proto

# Build the service
build:
	go build -o bin/accounts-service cmd/main.go

# Run the service
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy

# Run linter
lint:
	golangci-lint run 