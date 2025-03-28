# Create new directory structure
$directories = @(
    ".github/workflows",
    "api/openapi",
    "api/proto",
    "build/ci",
    "build/package",
    "config",
    "deployments/docker",
    "deployments/kubernetes/helm",
    "docs/architecture",
    "docs/api",
    "docs/development",
    "pkg/logger",
    "pkg/metrics",
    "pkg/validator",
    "scripts",
    "services"
)

foreach ($dir in $directories) {
    New-Item -ItemType Directory -Force -Path $dir
}

# Move existing services to services directory
$services = @(
    "aml-service",
    "auth-service",
    "investment-service",
    "kyc-service",
    "credit-scoring-service"
)

foreach ($service in $services) {
    if (Test-Path $service) {
        Move-Item -Path $service -Destination "services/"
    }
}

# Move Helm charts to deployments
if (Test-Path "helm") {
    Move-Item -Path "helm/*" -Destination "deployments/kubernetes/helm/"
}

# Move observability components to deployments
$observability = @("jaeger", "grafana", "prometheus")
foreach ($component in $observability) {
    if (Test-Path $component) {
        Move-Item -Path $component -Destination "deployments/kubernetes/"
    }
}

# Move configuration
if (Test-Path "config") {
    Move-Item -Path "config/*" -Destination "config/"
}

# Move scripts
if (Test-Path "scripts") {
    Move-Item -Path "scripts/*" -Destination "scripts/"
}

# Move GitHub Actions workflows
if (Test-Path ".github/workflows") {
    Move-Item -Path ".github/workflows/*" -Destination ".github/workflows/"
}

# Move development docker-compose
if (Test-Path "docker-compose.dev.yml") {
    Move-Item -Path "docker-compose.dev.yml" -Destination "deployments/docker/"
}

# Create service-specific directories for each service
foreach ($service in $services) {
    $serviceDirs = @(
        "cmd/server",
        "internal/domain",
        "internal/infrastructure",
        "internal/interfaces",
        "internal/usecase",
        "configs",
        "deployments",
        "docs",
        "test"
    )
    
    foreach ($dir in $serviceDirs) {
        New-Item -ItemType Directory -Force -Path "services/$service/$dir"
    }
}

# Create root go.mod if it doesn't exist
if (-not (Test-Path "go.mod")) {
    @"
module github.com/sparkfund

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v8 v8.11.5
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/google/uuid v1.3.0
    github.com/joho/godotenv v1.5.1
    github.com/prometheus/client_golang v1.16.0
    github.com/sirupsen/logrus v1.9.3
    gorm.io/driver/postgres v1.4.7
    gorm.io/gorm v1.25.5
)
"@ | Out-File -FilePath "go.mod" -Encoding UTF8
}

Write-Host "Project structure reorganized successfully!" 