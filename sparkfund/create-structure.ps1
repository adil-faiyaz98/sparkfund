$services = @(
    "api-gateway",
    "kyc-service",
    "aml-service",
    "fraud-detection-service",
    "credit-scoring-service",
    "risk-management-service",
    "notification-service",
    "consent-management-service",
    "logging-service",
    "security-service",
    "email-service",
    "blockchain-service"
)

$commonDirs = @(
    "cmd",
    "config",
    "internal/handlers",
    "internal/models",
    "internal/repositories",
    "internal/services",
    "internal/thirdparty",
    "internal/utils",
    "pkg"
)

foreach ($service in $services) {
    Write-Host "Creating structure for $service..."
    
    # Create main directories
    foreach ($dir in $commonDirs) {
        New-Item -ItemType Directory -Path "$service/$dir" -Force
    }
    
    # Create main.go in cmd directory
    $mainContent = @"
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/adil-faiyaz98/sparkfund/$service/config"
    "github.com/adil-faiyaz98/sparkfund/$service/internal/handlers"
    "github.com/adil-faiyaz98/sparkfund/$service/internal/services"
    "github.com/adil-faiyaz98/sparkfund/$service/internal/repositories"
    "github.com/adil-faiyaz98/sparkfund/$service/pkg/database"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // Initialize database connection
    dbConfig := database.NewConfig()
    db, err := database.NewConnection(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize repositories
    repo := repositories.NewRepository(db)

    // Initialize services
    svc := services.NewService(repo)

    // Initialize handlers
    handler := handlers.NewHandler(svc)

    // Set up router
    r := mux.NewRouter()
    handler.RegisterRoutes(r)

    // Start server
    log.Printf("Starting $service on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
"@
    Set-Content -Path "$service/cmd/main.go" -Value $mainContent

    # Create config.go
    $configContent = @"
package config

import (
    "os"
)

type Config struct {
    Port     string
    Database DatabaseConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Port: os.Getenv("PORT"),
        Database: DatabaseConfig{
            Host:     os.Getenv("DB_HOST"),
            Port:     os.Getenv("DB_PORT"),
            User:     os.Getenv("DB_USER"),
            Password: os.Getenv("DB_PASSWORD"),
            Name:     os.Getenv("DB_NAME"),
        },
    }
    return cfg, nil
}
"@
    Set-Content -Path "$service/config/config.go" -Value $configContent

    # Create go.mod
    $goModContent = @"
module github.com/adil-faiyaz98/sparkfund/$service

go 1.20

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
)
"@
    Set-Content -Path "$service/go.mod" -Value $goModContent

    # Create Dockerfile
    $dockerfileContent = @"
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
"@
    Set-Content -Path "$service/Dockerfile" -Value $dockerfileContent

    # Create .gitignore
    $gitignoreContent = @"
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
"@
    Set-Content -Path "$service/.gitignore" -Value $gitignoreContent
}

Write-Host "Directory structure created successfully!" 