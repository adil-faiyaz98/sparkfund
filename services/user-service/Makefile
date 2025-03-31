.PHONY: build run test migrate-up migrate-down clean

# Build the application
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/migrate cmd/migrate/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run migrations up
migrate-up:
	go run cmd/migrate/main.go -up

# Run migrations down
migrate-down:
	go run cmd/migrate/main.go -down

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Generate mocks (if using mockgen)
mocks:
	mockgen -source=internal/repository/user_repository.go -destination=internal/repository/mocks/user_repository_mock.go

# Run all checks
check: fmt lint test 