# KYC and AI Services

This project contains two microservices:
1. **KYC Service**: Handles KYC verification functionality
2. **AI Service**: Provides AI-powered capabilities for KYC verification

## Prerequisites

- Go (1.16 or later)
- Python (3.6 or later)
- Windows OS

## Files

- `simple_ai_server.py`: The AI service implementation
- `simple_kyc_server.go`: The KYC service implementation
- `run_kyc_ai_services.bat`: Script to start both services
- `test_kyc_ai_services.ps1`: PowerShell script to test the services

## Deployment

1. Clone this repository to your local machine
2. Open a command prompt in the repository directory
3. Run the deployment script:

```
run_kyc_ai_services.bat
```

This will start both services:
- AI Service on port 8001
- KYC Service on port 8081

## Testing

1. Open a PowerShell window in the repository directory
2. Run the test script:

```powershell
.\test_kyc_ai_services.ps1
```

This will test all the endpoints and display the results.

## API Endpoints

### KYC Service (http://localhost:8081)

- `GET /health`: Health check endpoint
- `POST /api/v1/auth/login`: Login endpoint
- `GET /api/v1/ai/models`: Get AI models
- `POST /api/v1/verifications`: Create a verification
- `GET /api/v1/verifications`: Get verifications
- `POST /api/v1/ai/analyze-document`: Analyze a document
- `POST /api/v1/ai/match-faces`: Match faces
- `POST /api/v1/ai/analyze-risk`: Analyze risk
- `POST /api/v1/ai/detect-anomalies`: Detect anomalies
- `POST /api/v1/ai/process-document`: Process a document

### AI Service (http://localhost:8001)

- `GET /health`: Health check endpoint
- `GET /api/v1/models`: Get AI models
- `POST /api/v1/document/analyze`: Analyze a document
- `POST /api/v1/document/analyze-base64`: Analyze a document from base64 encoded image
- `GET /api/v1/document/types`: Get document types
- `POST /api/v1/face/match`: Match faces
- `POST /api/v1/face/match-base64`: Match faces from base64 encoded images
- `GET /api/v1/face/thresholds`: Get face matching thresholds
- `POST /api/v1/risk/analyze`: Analyze risk
- `GET /api/v1/risk/factors`: Get risk factors
- `GET /api/v1/risk/levels`: Get risk levels
- `POST /api/v1/anomaly/detect`: Detect anomalies
- `GET /api/v1/anomaly/types`: Get anomaly types
- `POST /api/v1/process-document`: Process a document

## Stopping the Services

To stop the services, press any key in the command prompt window where you ran the deployment script.

## Manual Testing

You can also test the services manually using tools like Postman or curl. Here are some example requests:

### Login

```
POST http://localhost:8081/api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "password123",
  "device_info": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0",
    "device_type": "Desktop",
    "os": "Windows",
    "browser": "Chrome",
    "location": "New York, USA"
  }
}
```

### Create Verification

```
POST http://localhost:8081/api/v1/verifications
Content-Type: application/json

{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "kyc_id": "123e4567-e89b-12d3-a456-426614174002",
  "document_id": "123e4567-e89b-12d3-a456-426614174003",
  "method": "AI",
  "status": "PENDING"
}
```

### Analyze Document

```
POST http://localhost:8081/api/v1/ai/analyze-document
Content-Type: application/json

{
  "document_id": "123e4567-e89b-12d3-a456-426614174003",
  "verification_id": "123e4567-e89b-12d3-a456-426614174001"
}
```

### Match Faces

```
POST http://localhost:8081/api/v1/ai/match-faces
Content-Type: application/json

{
  "document_id": "123e4567-e89b-12d3-a456-426614174003",
  "selfie_id": "123e4567-e89b-12d3-a456-426614174004",
  "verification_id": "123e4567-e89b-12d3-a456-426614174001"
}
```

### Analyze Risk

```
POST http://localhost:8081/api/v1/ai/analyze-risk
Content-Type: application/json

{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "verification_id": "123e4567-e89b-12d3-a456-426614174001",
  "device_info": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0",
    "device_type": "Desktop",
    "os": "Windows",
    "browser": "Chrome",
    "location": "New York, USA"
  }
}
```

### Detect Anomalies

```
POST http://localhost:8081/api/v1/ai/detect-anomalies
Content-Type: application/json

{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "verification_id": "123e4567-e89b-12d3-a456-426614174001",
  "device_info": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0",
    "device_type": "Desktop",
    "os": "Windows",
    "browser": "Chrome",
    "location": "New York, USA"
  }
}
```
