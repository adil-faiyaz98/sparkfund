# Function to check if a URL is accessible
function Test-Url {
    param($url)
    try {
        $response = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 5
        return $response.StatusCode -eq 200
    }
    catch {
        return $false
    }
}

# Function to wait for a service to be ready
function Wait-ForService {
    param($url, $serviceName)
    Write-Host "Waiting for $serviceName to be ready..."
    $maxAttempts = 30
    $attempt = 0
    while ($attempt -lt $maxAttempts) {
        if (Test-Url $url) {
            Write-Host "$serviceName is ready!"
            return $true
        }
        $attempt++
        Start-Sleep -Seconds 2
    }
    Write-Host "Timeout waiting for $serviceName"
    return $false
}

# List of endpoints to check and open
$endpoints = @(
    @{
        Name = "Investment Management Interface"
        Url = "http://localhost:8081/api/v1/investments/"
    },
    @{
        Name = "API Gateway Health"
        Url = "http://localhost:8080/health"
    },
    @{
        Name = "API Gateway Metrics"
        Url = "http://localhost:8080/metrics"
    },
    @{
        Name = "Investment Service Health"
        Url = "http://localhost:8081/health"
    },
    @{
        Name = "Investment Service Metrics"
        Url = "http://localhost:8081/metrics"
    },
    @{
        Name = "Prometheus"
        Url = "http://localhost:9090"
    },
    @{
        Name = "Grafana"
        Url = "http://localhost:3000"
    }
)

# Wait for services to be ready
Write-Host "Checking if services are ready..."
foreach ($endpoint in $endpoints) {
    Wait-ForService -url $endpoint.Url -serviceName $endpoint.Name
}

# Open Chrome tabs for each endpoint
Write-Host "Opening endpoints in Chrome..."
$chromePath = "C:\Program Files\Google\Chrome\Application\chrome.exe"
$urls = $endpoints.Url -join " "
Start-Process -FilePath $chromePath -ArgumentList $urls

# Generate and display JWT token
Write-Host "`nGenerating JWT token..."
cd services/api-gateway/scripts
$token = go run generate_token.go
Write-Host "JWT Token: $token"

# Run seed script
Write-Host "`nRunning database seed script..."
cd ../../investment-service/scripts
go run seed.go

Write-Host "`nAll endpoints are open in Chrome tabs!"
Write-Host "You can now test the API endpoints using the generated JWT token." 