# Service Template

This is a template for creating new services in the SparkFund platform.

## Directory Structure

```
service-template/
├── cmd/                    # Command-line applications
│   └── main.go             # Main entry point
├── config/                 # Configuration files
│   ├── config.base.yaml    # Base configuration
│   ├── config.dev.yaml     # Development configuration
│   └── config.prod.yaml    # Production configuration
├── docs/                   # Documentation
│   └── swagger/            # Swagger documentation
├── internal/               # Private application code
│   ├── api/                # API layer
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # HTTP middleware
│   │   └── routes/         # HTTP routes
│   ├── config/             # Configuration code
│   ├── database/           # Database code
│   │   └── migrations/     # Database migrations
│   ├── domain/             # Domain layer
│   │   ├── model/          # Domain models
│   │   └── repository/     # Repository interfaces
│   ├── server/             # Server code
│   └── service/            # Service layer
├── pkg/                    # Public code that can be imported by other services
├── scripts/                # Scripts for development, CI/CD, etc.
├── test/                   # Test files
│   ├── integration/        # Integration tests
│   └── mocks/              # Mock implementations
├── .air.toml               # Air configuration for hot reloading
├── .dockerignore           # Docker ignore file
├── .gitignore              # Git ignore file
├── Dockerfile              # Dockerfile for building the service
├── Makefile                # Makefile for common tasks
├── docker-compose.yml      # Docker Compose for production
├── docker-compose.dev.yml  # Docker Compose for development
├── go.mod                  # Go module file
└── go.sum                  # Go module checksum
```

## Getting Started

To create a new service based on this template:

1. Copy the template directory:

```bash
cp -r services/service-template services/your-service-name
```

2. Update the service name in all files:

```bash
cd services/your-service-name
find . -type f -exec sed -i 's/service-template/your-service-name/g' {} \;
```

3. Update the module name in go.mod:

```bash
sed -i 's|github.com/adil-faiyaz98/sparkfund/services/service-template|github.com/adil-faiyaz98/sparkfund/services/your-service-name|g' go.mod
```

4. Initialize the Git repository:

```bash
git add .
git commit -m "Create new service: your-service-name"
```

5. Create the Kubernetes deployment file:

```bash
cp deploy/k8s/templates/service-template.yaml deploy/k8s/your-service-name.yaml
sed -i 's/SERVICE_NAME/your-service-name/g' deploy/k8s/your-service-name.yaml
```

6. Create the GitHub Actions workflow:

```bash
cp .github/workflows/service-template.yml .github/workflows/your-service-name.yml
sed -i 's/SERVICE_NAME/your-service-name/g' .github/workflows/your-service-name.yml
```

7. Update the README.md with your service-specific information.

## Development

### Running the Service

```bash
# Using Go
go run cmd/main.go

# Using Air (hot reloading)
air

# Using Docker Compose
docker-compose -f docker-compose.dev.yml up
```

### Running Tests

```bash
go test ./...
```

### Building the Service

```bash
go build -o build/service cmd/main.go
```

### Building the Docker Image

```bash
docker build -t your-service-name .
```

## Deployment

### Kubernetes

```bash
kubectl apply -f deploy/k8s/your-service-name.yaml
```

### CI/CD

The service is automatically built, tested, and deployed using GitHub Actions. See `.github/workflows/your-service-name.yml` for details.
