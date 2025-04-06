# Test script for KYC and AI services

$kycServiceUrl = "http://localhost:8081"
$aiServiceUrl = "http://localhost:8001"

Write-Host "Testing KYC and AI services..." -ForegroundColor Green
Write-Host ""

# Test health endpoints
Write-Host "Testing health endpoints..." -ForegroundColor Yellow
try {
    $kycHealth = Invoke-RestMethod -Uri "$kycServiceUrl/health" -Method Get
    Write-Host "KYC Service health: $($kycHealth.status)" -ForegroundColor Green
} catch {
    Write-Host "Error testing KYC Service health: $_" -ForegroundColor Red
}

try {
    $aiHealth = Invoke-RestMethod -Uri "$aiServiceUrl/health" -Method Get
    Write-Host "AI Service health: $($aiHealth.status)" -ForegroundColor Green
} catch {
    Write-Host "Error testing AI Service health: $_" -ForegroundColor Red
}
Write-Host ""

# Test login
Write-Host "Testing login..." -ForegroundColor Yellow
$loginBody = @{
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

try {
    $loginResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    Write-Host "Login successful. Token: $($loginResponse.token.Substring(0, 20))..." -ForegroundColor Green
    $token = $loginResponse.token
    $userId = $loginResponse.user.id
} catch {
    Write-Host "Error testing login: $_" -ForegroundColor Red
    $userId = "123e4567-e89b-12d3-a456-426614174000"
}
Write-Host ""

# Test AI models
Write-Host "Testing AI models..." -ForegroundColor Yellow
try {
    $modelsResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/ai/models" -Method Get
    Write-Host "AI models retrieved successfully. Count: $($modelsResponse.models.Count)" -ForegroundColor Green
    foreach ($model in $modelsResponse.models) {
        Write-Host "  - $($model.name) (Type: $($model.type), Accuracy: $($model.accuracy))" -ForegroundColor Cyan
    }
} catch {
    Write-Host "Error testing AI models: $_" -ForegroundColor Red
}
Write-Host ""

# Create verification
Write-Host "Creating verification..." -ForegroundColor Yellow
$verificationBody = @{
    user_id = $userId
    kyc_id = "123e4567-e89b-12d3-a456-426614174002"
    document_id = "123e4567-e89b-12d3-a456-426614174003"
    method = "AI"
    status = "PENDING"
} | ConvertTo-Json

try {
    $verificationResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/verifications" -Method Post -Body $verificationBody -ContentType "application/json"
    Write-Host "Verification created successfully. ID: $($verificationResponse.verification.id)" -ForegroundColor Green
    $verificationId = $verificationResponse.verification.id
    $documentId = $verificationResponse.verification.document_id
} catch {
    Write-Host "Error creating verification: $_" -ForegroundColor Red
    $verificationId = "123e4567-e89b-12d3-a456-426614174001"
    $documentId = "123e4567-e89b-12d3-a456-426614174003"
}
Write-Host ""

# Test document analysis
Write-Host "Testing document analysis..." -ForegroundColor Yellow
$documentBody = @{
    document_id = $documentId
    verification_id = $verificationId
} | ConvertTo-Json

try {
    $documentResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/ai/analyze-document" -Method Post -Body $documentBody -ContentType "application/json"
    Write-Host "Document analysis successful." -ForegroundColor Green
    Write-Host "  - Document type: $($documentResponse.document_type)" -ForegroundColor Cyan
    Write-Host "  - Is authentic: $($documentResponse.is_authentic)" -ForegroundColor Cyan
    Write-Host "  - Confidence: $($documentResponse.confidence)" -ForegroundColor Cyan
    Write-Host "  - Extracted data: $($documentResponse.extracted_data | ConvertTo-Json -Compress)" -ForegroundColor Cyan
} catch {
    Write-Host "Error testing document analysis: $_" -ForegroundColor Red
}
Write-Host ""

# Test face matching
Write-Host "Testing face matching..." -ForegroundColor Yellow
$faceBody = @{
    document_id = $documentId
    selfie_id = "123e4567-e89b-12d3-a456-426614174004"
    verification_id = $verificationId
} | ConvertTo-Json

try {
    $faceResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/ai/match-faces" -Method Post -Body $faceBody -ContentType "application/json"
    Write-Host "Face matching successful." -ForegroundColor Green
    Write-Host "  - Is match: $($faceResponse.is_match)" -ForegroundColor Cyan
    Write-Host "  - Confidence: $($faceResponse.confidence)" -ForegroundColor Cyan
} catch {
    Write-Host "Error testing face matching: $_" -ForegroundColor Red
}
Write-Host ""

# Test risk analysis
Write-Host "Testing risk analysis..." -ForegroundColor Yellow
$riskBody = @{
    user_id = $userId
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

try {
    $riskResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/ai/analyze-risk" -Method Post -Body $riskBody -ContentType "application/json"
    Write-Host "Risk analysis successful." -ForegroundColor Green
    Write-Host "  - Risk score: $($riskResponse.risk_score)" -ForegroundColor Cyan
    Write-Host "  - Risk level: $($riskResponse.risk_level)" -ForegroundColor Cyan
    Write-Host "  - Risk factors: $($riskResponse.risk_factors -join ', ')" -ForegroundColor Cyan
} catch {
    Write-Host "Error testing risk analysis: $_" -ForegroundColor Red
}
Write-Host ""

# Test anomaly detection
Write-Host "Testing anomaly detection..." -ForegroundColor Yellow
$anomalyBody = @{
    user_id = $userId
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

try {
    $anomalyResponse = Invoke-RestMethod -Uri "$kycServiceUrl/api/v1/ai/detect-anomalies" -Method Post -Body $anomalyBody -ContentType "application/json"
    Write-Host "Anomaly detection successful." -ForegroundColor Green
    Write-Host "  - Is anomaly: $($anomalyResponse.is_anomaly)" -ForegroundColor Cyan
    Write-Host "  - Anomaly score: $($anomalyResponse.anomaly_score)" -ForegroundColor Cyan
    if ($anomalyResponse.is_anomaly) {
        Write-Host "  - Anomaly type: $($anomalyResponse.anomaly_type)" -ForegroundColor Cyan
        Write-Host "  - Reasons: $($anomalyResponse.reasons -join ', ')" -ForegroundColor Cyan
    }
} catch {
    Write-Host "Error testing anomaly detection: $_" -ForegroundColor Red
}
Write-Host ""

Write-Host "All tests completed!" -ForegroundColor Green
Write-Host ""
Write-Host "Summary:" -ForegroundColor Yellow
Write-Host "- KYC Service: http://localhost:8081" -ForegroundColor Cyan
Write-Host "- AI Service: http://localhost:8001" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press any key to exit..." -ForegroundColor Green
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
