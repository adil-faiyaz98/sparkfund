# AI Service for KYC

This service provides AI-powered capabilities for KYC (Know Your Customer) verification, including document verification, facial recognition, risk analysis, and anomaly detection.

## Features

- Document verification using computer vision and OCR
- Facial recognition for identity verification
- Risk analysis based on user data and device information
- Anomaly detection to identify suspicious behavior
- AI model management and versioning

## Prerequisites

- Docker and Docker Compose
- Python 3.9 or later (for local development)

## Getting Started

### Running with Docker Compose

1. Clone the repository
2. Navigate to the `services/ai-service` directory
3. Run the service using Docker Compose:

```bash
docker-compose up --build
```

The service will be available at http://localhost:8000.

### Running Locally

1. Clone the repository
2. Navigate to the `services/ai-service` directory
3. Create a virtual environment:

```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

4. Install dependencies:

```bash
pip install -r requirements.txt
```

5. Run the service:

```bash
uvicorn src.main:app --reload
```

The service will be available at http://localhost:8000.

## API Documentation

The API documentation is available at http://localhost:8000/docs when the service is running.

## API Endpoints

### Document Verification

- `POST /api/v1/document/analyze` - Analyze a document for authenticity
- `POST /api/v1/document/analyze-base64` - Analyze a document from base64 encoded image
- `GET /api/v1/document/types` - Get a list of supported document types

### Face Recognition

- `POST /api/v1/face/match` - Match a selfie with a document photo
- `POST /api/v1/face/match-base64` - Match faces from base64 encoded images
- `GET /api/v1/face/thresholds` - Get face matching thresholds

### Risk Analysis

- `POST /api/v1/risk/analyze` - Analyze risk based on user data
- `GET /api/v1/risk/factors` - Get a list of risk factors
- `GET /api/v1/risk/levels` - Get risk level thresholds

### Anomaly Detection

- `POST /api/v1/anomaly/detect` - Detect anomalies in user behavior
- `GET /api/v1/anomaly/types` - Get a list of anomaly types

### AI Models

- `GET /api/v1/models` - List all AI models
- `GET /api/v1/models/{model_id}` - Get AI model information
- `GET /api/v1/models/type/{model_type}` - Get latest AI model by type
- `POST /api/v1/models/upload/{model_type}` - Upload a new AI model

## Directory Structure

- `src/`: Application source code
  - `main.py`: Main application entry point
  - `config.py`: Configuration settings
  - `models/`: Data models
  - `routers/`: API routers
  - `services/`: Business logic services
- `models/`: AI models directory
  - `document/`: Document verification models
  - `face/`: Face recognition models
  - `risk/`: Risk analysis models
  - `anomaly/`: Anomaly detection models
- `tests/`: Test files
- `uploads/`: Temporary upload directory

## Environment Variables

The service can be configured using the following environment variables:

- `DEBUG`: Enable debug mode (default: False)
- `HOST`: Host to bind to (default: 0.0.0.0)
- `PORT`: Port to listen on (default: 8000)
- `API_KEY`: API key for authentication (default: your-api-key)
- `MODEL_PATH`: Path to AI models directory (default: ./models)
- `UPLOAD_DIR`: Path to uploads directory (default: ./uploads)

## Integration with KYC Service

The AI service is designed to be used by the KYC service. The KYC service makes API calls to the AI service to perform document verification, facial recognition, risk analysis, and anomaly detection.

To integrate with the KYC service, set the following environment variables in the KYC service:

- `AI_SERVICE_URL`: URL of the AI service (e.g., http://ai-service:8000)
- `AI_SERVICE_API_KEY`: API key for the AI service

## License

This project is licensed under the MIT License - see the LICENSE file for details.
