# PowerShell script to test all available endpoints in SparkFund

Write-Host "====================================================="
Write-Host "SparkFund - Testing Available Endpoints"
Write-Host "====================================================="
Write-Host ""

# Define the JWT token for authentication
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxOTE2MjM5MDIyLCJyb2xlcyI6WyJhZG1pbiIsInVzZXIiXX0.Ks0I-dCdjWUxJEwuGP0qlyYJGXXjUYlLCRwPIZXI5Ss"
$headers = @{
    "Content-Type" = "application/json"
    "Authorization" = "Bearer $token"
}

# Define endpoints to test
$endpoints = @(
    # API Gateway
    @{Name = "API Gateway"; Url = "http://localhost:8080/health"; Method = "GET"; Body = $null},
    
    # KYC Service
    @{Name = "KYC Health"; Url = "http://localhost:8080/api/kyc/health"; Method = "GET"; Body = $null},
    @{Name = "KYC Status"; Url = "http://localhost:8080/api/kyc/api/v1/kyc/status"; Method = "GET"; Body = $null},
    @{Name = "KYC Verify"; Url = "http://localhost:8080/api/kyc/api/v1/kyc/verify"; Method = "POST"; Body = "{}"},
    
    # Investment Service
    @{Name = "Investment Health"; Url = "http://localhost:8080/api/investment/health"; Method = "GET"; Body = $null},
    @{Name = "Investment List"; Url = "http://localhost:8080/api/investment/api/v1/investments"; Method = "GET"; Body = $null},
    @{Name = "Investment Create"; Url = "http://localhost:8080/api/investment/api/v1/investments/create"; Method = "POST"; Body = '{"user_id": "123", "amount": 1000, "type": "STOCK", "symbol": "AAPL", "quantity": 10}'},
    
    # User Service
    @{Name = "User Health"; Url = "http://localhost:8080/api/user/health"; Method = "GET"; Body = $null},
    @{Name = "User List"; Url = "http://localhost:8080/api/user/api/v1/users"; Method = "GET"; Body = $null},
    @{Name = "User Register"; Url = "http://localhost:8080/api/user/api/v1/users/register"; Method = "POST"; Body = '{"email": "test@example.com", "first_name": "Test", "last_name": "User", "password": "password123"}'},
    
    # AI Service
    @{Name = "AI Health"; Url = "http://localhost:8080/api/ai/health"; Method = "GET"; Body = $null}
)

# Test each endpoint
foreach ($endpoint in $endpoints) {
    Write-Host "Testing $($endpoint.Name): $($endpoint.Url)" -NoNewline
    
    try {
        if ($endpoint.Method -eq "GET") {
            $response = Invoke-WebRequest -Uri $endpoint.Url -Method $endpoint.Method -Headers $headers -TimeoutSec 5 -ErrorAction Stop
        } else {
            $response = Invoke-WebRequest -Uri $endpoint.Url -Method $endpoint.Method -Headers $headers -Body $endpoint.Body -TimeoutSec 5 -ErrorAction Stop
        }
        
        if ($response.StatusCode -eq 200) {
            Write-Host " [OK]" -ForegroundColor Green
            Write-Host "  Response: $($response.Content)" -ForegroundColor Gray
        } else {
            Write-Host " [FAIL] Status code: $($response.StatusCode)" -ForegroundColor Red
        }
    } catch {
        Write-Host " [FAIL] Error: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    Write-Host ""
}

Write-Host "====================================================="
Write-Host "Testing Complete"
Write-Host "====================================================="
