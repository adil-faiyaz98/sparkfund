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

foreach ($service in $services) {
    Write-Host "Processing $service..."
    
    # Delete main.go if it exists
    if (Test-Path "$service/main.go") {
        Remove-Item "$service/main.go" -Force
        Write-Host "Deleted main.go from $service"
    }
    
    # Ensure go.mod exists
    if (-not (Test-Path "$service/go.mod")) {
        $goModContent = @"
module github.com/adil-faiyaz98/sparkfund/$service

go 1.20

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
)
"@
        Set-Content -Path "$service/go.mod" -Value $goModContent
        Write-Host "Created go.mod for $service"
    }
    
    # Run go mod tidy to generate go.sum
    Push-Location $service
    go mod tidy
    Pop-Location
}

Write-Host "Cleanup completed successfully!" 