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
    Write-Host "Updating $service..."
    $goModPath = Join-Path $service "go.mod"
    $content = Get-Content $goModPath
    $content[0] = "module github.com/adil-faiyaz98/sparkfund/$service"
    $content | Set-Content $goModPath
}

Write-Host "All module paths have been updated successfully!" 