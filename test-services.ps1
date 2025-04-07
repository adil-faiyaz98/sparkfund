# SparkFund - Test Services Script
Write-Host "====================================================="
Write-Host "SparkFund - Testing Services"
Write-Host "====================================================="

# Function to test an endpoint
function Test-Endpoint {
    param (
        [string]$Service,
        [string]$Endpoint,
        [string]$Method = "GET"
    )
    
    Write-Host "Testing $Service at $Endpoint..."
    
    try {
        $response = Invoke-WebRequest -Uri $Endpoint -Method $Method -ErrorAction Stop
        Write-Host "Status: $($response.StatusCode) - Success!" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "Status: Failed - $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Test API Gateway
$apiGatewayHealth = Test-Endpoint -Service "API Gateway" -Endpoint "http://localhost:8080/health"

# Test KYC Service
$kycServiceHealth = Test-Endpoint -Service "KYC Service" -Endpoint "http://localhost:8081/health"

# Test Investment Service
$investmentServiceHealth = Test-Endpoint -Service "Investment Service" -Endpoint "http://localhost:8082/health"

# Test User Service
$userServiceHealth = Test-Endpoint -Service "User Service" -Endpoint "http://localhost:8083/health"

# Test AI Service
$aiServiceHealth = Test-Endpoint -Service "AI Service" -Endpoint "http://localhost:8001/health"

# Summary
Write-Host "====================================================="
Write-Host "Service Health Summary:"
Write-Host "====================================================="
Write-Host "API Gateway: $(if ($apiGatewayHealth) { "Healthy" } else { "Unhealthy" })"
Write-Host "KYC Service: $(if ($kycServiceHealth) { "Healthy" } else { "Unhealthy" })"
Write-Host "Investment Service: $(if ($investmentServiceHealth) { "Healthy" } else { "Unhealthy" })"
Write-Host "User Service: $(if ($userServiceHealth) { "Healthy" } else { "Unhealthy" })"
Write-Host "AI Service: $(if ($aiServiceHealth) { "Healthy" } else { "Unhealthy" })"
Write-Host "====================================================="

# Test basic functionality if all services are healthy
if ($apiGatewayHealth -and $kycServiceHealth -and $investmentServiceHealth -and $userServiceHealth -and $aiServiceHealth) {
    Write-Host "All services are healthy. Testing basic functionality..."
    
    # Add more functional tests here
    
    Write-Host "Basic functionality tests completed."
}
else {
    Write-Host "Some services are unhealthy. Skipping functional tests." -ForegroundColor Yellow
}

Write-Host "====================================================="
