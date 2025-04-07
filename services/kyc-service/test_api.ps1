# Base URL
$BASE_URL = "http://localhost:8080/api/v1"

# Function to make API requests
function Make-Request {
    param (
        [string]$Method,
        [string]$Endpoint,
        [string]$Data,
        [string]$Token
    )
    
    Write-Host "Making $Method request to $Endpoint" -ForegroundColor Green
    
    if ($Data) {
        Write-Host "Request data: $Data"
    }
    
    $headers = @{
        "Content-Type" = "application/json"
    }
    
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    
    $params = @{
        Method = $Method
        Uri = "$BASE_URL$Endpoint"
        Headers = $headers
    }
    
    if ($Data) {
        $params["Body"] = $Data
    }
    
    try {
        $response = Invoke-RestMethod @params
        Write-Host "Response: $($response | ConvertTo-Json -Depth 10)"
        Write-Host ""
        return $response
    }
    catch {
        Write-Host "Error: $_" -ForegroundColor Red
        Write-Host "Status Code: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

# Test health endpoint
Write-Host "Testing health endpoint" -ForegroundColor Green
Invoke-RestMethod -Uri "http://localhost:8080/health"
Write-Host ""

# Login
Write-Host "Testing login" -ForegroundColor Green
$loginData = @{
    email = "admin@example.com"
    password = "password123"
    device_info = @{
        ip_address = "192.168.1.1"
        user_agent = "Mozilla/5.0"
        device_type = "Desktop"
        os = "Windows"
        browser = "Chrome"
        location = "New York, USA"
    }
} | ConvertTo-Json

$loginResponse = Make-Request -Method "POST" -Endpoint "/auth/login" -Data $loginData

# Extract token from login response
$token = $loginResponse.token
Write-Host "Token: $token" -ForegroundColor Green

# Test AI models endpoint
Write-Host "Testing AI models endpoint" -ForegroundColor Green
Make-Request -Method "GET" -Endpoint "/ai/models" -Token $token

# Upload a document
Write-Host "Testing document upload" -ForegroundColor Green
$documentForm = @{
    file = Get-Item "./test_data/passport.jpg"
    type = "PASSPORT"
    name = "passport.jpg"
}

$documentResponse = Invoke-RestMethod -Method POST -Uri "$BASE_URL/documents" -Headers @{
    "Authorization" = "Bearer $token"
} -Form $documentForm

Write-Host "Response: $($documentResponse | ConvertTo-Json -Depth 10)"
Write-Host ""

# Extract document ID
$documentId = $documentResponse.id
Write-Host "Document ID: $documentId" -ForegroundColor Green

# Upload a selfie
Write-Host "Testing selfie upload" -ForegroundColor Green
$selfieForm = @{
    file = Get-Item "./test_data/selfie.jpg"
    type = "SELFIE"
    name = "selfie.jpg"
}

$selfieResponse = Invoke-RestMethod -Method POST -Uri "$BASE_URL/documents" -Headers @{
    "Authorization" = "Bearer $token"
} -Form $selfieForm

Write-Host "Response: $($selfieResponse | ConvertTo-Json -Depth 10)"
Write-Host ""

# Extract selfie ID
$selfieId = $selfieResponse.id
Write-Host "Selfie ID: $selfieId" -ForegroundColor Green

# Create a verification
Write-Host "Testing verification creation" -ForegroundColor Green
$verificationData = @{
    user_id = $loginResponse.user.id
    kyc_id = [guid]::NewGuid().ToString()
    document_id = $documentId
    method = "AI"
    status = "PENDING"
} | ConvertTo-Json

$verificationResponse = Make-Request -Method "POST" -Endpoint "/verifications" -Data $verificationData -Token $token

# Extract verification ID
$verificationId = $verificationResponse.verification.id
Write-Host "Verification ID: $verificationId" -ForegroundColor Green

# Test document analysis
Write-Host "Testing document analysis" -ForegroundColor Green
$analysisData = @{
    document_id = $documentId
    verification_id = $verificationId
} | ConvertTo-Json

Make-Request -Method "POST" -Endpoint "/ai/analyze-document" -Data $analysisData -Token $token

# Test face matching
Write-Host "Testing face matching" -ForegroundColor Green
$matchData = @{
    document_id = $documentId
    selfie_id = $selfieId
    verification_id = $verificationId
} | ConvertTo-Json

Make-Request -Method "POST" -Endpoint "/ai/match-faces" -Data $matchData -Token $token

# Test risk analysis
Write-Host "Testing risk analysis" -ForegroundColor Green
$riskData = @{
    user_id = $loginResponse.user.id
    verification_id = $verificationId
    device_info = @{
        ip_address = "192.168.1.1"
        user_agent = "Mozilla/5.0"
        device_type = "Desktop"
        os = "Windows"
        browser = "Chrome"
        location = "New York, USA"
    }
} | ConvertTo-Json

Make-Request -Method "POST" -Endpoint "/ai/analyze-risk" -Data $riskData -Token $token

# Test anomaly detection
Write-Host "Testing anomaly detection" -ForegroundColor Green
$anomalyData = @{
    user_id = $loginResponse.user.id
    verification_id = $verificationId
    device_info = @{
        ip_address = "192.168.1.1"
        user_agent = "Mozilla/5.0"
        device_type = "Desktop"
        os = "Windows"
        browser = "Chrome"
        location = "New York, USA"
    }
} | ConvertTo-Json

Make-Request -Method "POST" -Endpoint "/ai/detect-anomalies" -Data $anomalyData -Token $token

# Test document processing
Write-Host "Testing document processing" -ForegroundColor Green
$processData = @{
    document_id = $documentId
    selfie_id = $selfieId
    verification_id = $verificationId
    device_info = @{
        ip_address = "192.168.1.1"
        user_agent = "Mozilla/5.0"
        device_type = "Desktop"
        os = "Windows"
        browser = "Chrome"
        location = "New York, USA"
    }
} | ConvertTo-Json

Make-Request -Method "POST" -Endpoint "/ai/process-document" -Data $processData -Token $token

# Get verification status
Write-Host "Testing verification status" -ForegroundColor Green
Make-Request -Method "GET" -Endpoint "/verifications/$verificationId" -Token $token

Write-Host "All tests completed" -ForegroundColor Green
