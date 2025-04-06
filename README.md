# SparkFund Investment Platform

A microservices-based investment platform built with Go, featuring an API Gateway, Investment Service, KYC Service, and AI Service.

## Project Structure

```
.
├── services/
│   ├── api-gateway/        # API Gateway service
│   ├── investment-service/ # Investment Service
│   ├── kyc-service/        # KYC verification service
│   └── ai-service/         # AI service for document verification, facial recognition, etc.
├── deploy/
│   └── k8s/               # Kubernetes deployment configurations
└── docker-compose.yml     # Local development setup
```

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL
- Redis
- Make (optional, for using Makefile commands)

## Local Development

1. Clone the repository:

```bash
git clone https://github.com/adil-faiyaz98/sparkfund.git
cd sparkfund
```

2. Start the local development environment:

```bash
docker-compose up --build
```

This will start all required services:

- API Gateway (http://localhost:8080)
- Investment Service (http://localhost:8081)
- PostgreSQL (localhost:5432)
- Redis (localhost:6379)
- Prometheus (http://localhost:9090)
- Grafana (http://localhost:3000)

3. Alternatively, to run just the KYC and AI services:

```bash
.\run_services.bat
```

This will start:

- KYC Service (http://localhost:8081)
- AI Service (http://localhost:8001)

Swagger documentation is available at:

- KYC Service: http://localhost:8081/swagger-ui/
- AI Service: http://localhost:8001/docs

## API Endpoints

### Investment Service

- `POST /api/v1/investments` - Create a new investment
- `GET /api/v1/investments/:id` - Get investment by ID
- `GET /api/v1/investments` - List all investments
- `PUT /api/v1/investments/:id` - Update investment
- `DELETE /api/v1/investments/:id` - Delete investment

### Portfolio Management

- `POST /api/v1/portfolios` - Create a new portfolio
- `GET /api/v1/portfolios/:id` - Get portfolio by ID
- `PUT /api/v1/portfolios/:id` - Update portfolio
- `DELETE /api/v1/portfolios/:id` - Delete portfolio

### Transactions

- `POST /api/v1/transactions` - Create a new transaction

### KYC Service

- `GET /health` - Health check endpoint
- `POST /api/v1/auth/login` - Login endpoint
- `GET /api/v1/get-api-key` - Get API key for AI service
- `POST /api/v1/verifications` - Create a verification
- `GET /api/v1/verifications` - List all verifications
- `GET /api/v1/ai/models` - Get AI models
- `POST /api/v1/ai/analyze-document` - Analyze a document
- `POST /api/v1/ai/match-faces` - Match faces
- `POST /api/v1/ai/analyze-risk` - Analyze risk
- `POST /api/v1/ai/detect-anomalies` - Detect anomalies
- `POST /api/v1/ai/process-document` - Process a document

### AI Service

- `GET /health` - Health check endpoint
- `GET /api/v1/get-api-key` - Get API key for testing
- `POST /api/v1/document/analyze` - Analyze a document (file upload)
- `POST /api/v1/document/analyze-base64` - Analyze a document from base64 encoded image
- `GET /api/v1/document/types` - Get supported document types
- `POST /api/v1/face/match` - Match a selfie with a document photo (file upload)
- `POST /api/v1/face/match-base64` - Match faces from base64 encoded images
- `GET /api/v1/face/thresholds` - Get face matching thresholds
- `POST /api/v1/risk/analyze` - Analyze risk based on user data
- `GET /api/v1/risk/factors` - Get a list of risk factors
- `GET /api/v1/risk/levels` - Get risk level thresholds
- `POST /api/v1/anomaly/detect` - Detect anomalies in user behavior
- `GET /api/v1/anomaly/types` - Get a list of anomaly types
- `GET /api/v1/models` - List all AI models
- `GET /api/v1/models/{model_id}` - Get AI model information
- `GET /api/v1/models/type/{model_type}` - Get latest AI model by type

## Testing

### Running Tests

```bash
# Run tests for all services
make test

# Run tests for specific service
cd services/api-gateway && go test ./...
cd services/investment-service && go test ./...
```

### API Testing

You can use curl or any API testing tool like Postman to test the endpoints. Here's an example:

```bash
# Create a new investment
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123",
    "amount": 1000,
    "type": "STOCK",
    "symbol": "AAPL",
    "quantity": 10
  }'

# Get all investments
curl http://localhost:8080/api/v1/investments
```

## Monitoring

- Prometheus metrics: http://localhost:9090
- Grafana dashboards: http://localhost:3000 (username: admin, password: admin)

## Development Guidelines

1. Follow Go best practices and coding standards
2. Write tests for new features
3. Update documentation when making changes
4. Use meaningful commit messages

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
