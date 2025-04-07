# PowerShell script to check if all SparkFund services are running properly

Write-Host "====================================================="
Write-Host "SparkFund - Service Health Check"
Write-Host "====================================================="
Write-Host ""

# Check Docker container status
Write-Host "Checking Docker container status..."
docker-compose ps
Write-Host ""

# Define services to check
$services = @(
    @{Name = "API Gateway"; Url = "http://localhost:8080/health" },
    @{Name = "KYC Service"; Url = "http://localhost:8081/health" },
    @{Name = "Investment Service"; Url = "http://localhost:8082/health" },
    @{Name = "User Service"; Url = "http://localhost:8083/health" },
    @{Name = "AI Service"; Url = "http://localhost:8001/health" }
)

# Check each service
foreach ($service in $services) {
    Write-Host "Checking $($service.Name)..." -NoNewline

    try {
        $response = Invoke-WebRequest -Uri $service.Url -Method GET -TimeoutSec 5 -ErrorAction Stop

        if ($response.StatusCode -eq 200) {
            Write-Host " [OK]" -ForegroundColor Green
        }
        else {
            Write-Host " [FAIL] Status code: $($response.StatusCode)" -ForegroundColor Red
        }
    }
    catch {
        Write-Host " [FAIL] Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "====================================================="
Write-Host "Service URLs:"
Write-Host "====================================================="
Write-Host "API Gateway:         http://localhost:8080"
Write-Host "KYC Service:         http://localhost:8081"
Write-Host "Investment Service:  http://localhost:8082"
Write-Host "User Service:        http://localhost:8083"
Write-Host "AI Service:          http://localhost:8001"
Write-Host ""
Write-Host "Swagger Documentation:"
Write-Host "KYC Service:         http://localhost:8081/swagger-ui.html"
Write-Host "Investment Service:  http://localhost:8082/swagger-ui.html"
Write-Host "User Service:        http://localhost:8083/swagger-ui.html"
Write-Host "AI Service:          http://localhost:8001/docs"
Write-Host ""
Write-Host "Monitoring:"
Write-Host "Prometheus:          http://localhost:9090"
Write-Host "Grafana:             http://localhost:3000 (admin/admin)"
Write-Host "Jaeger:              http://localhost:16686"
Write-Host "====================================================="
