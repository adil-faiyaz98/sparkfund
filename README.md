# SparkFund Platform

<p align="center">
  <img src="https://via.placeholder.com/200x200?text=SparkFund" alt="SparkFund Logo" width="200" height="200">
</p>

SparkFund is a comprehensive financial platform with AI-powered investment recommendations, KYC verification, and user management capabilities. Built with a microservices architecture, it provides a scalable and robust foundation for financial applications.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Services](#-services)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [AI Capabilities](#-ai-capabilities)
- [Testing](#-testing)
- [Monitoring and Observability](#-monitoring-and-observability)
- [Deployment](#-deployment)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- **Microservices Architecture**: Scalable, maintainable, and resilient system design
- **AI-Powered Investment Recommendations**: Advanced algorithms for personalized investment advice
- **KYC Verification**: Secure and compliant customer onboarding with AI-powered document verification
- **User Management**: Comprehensive user authentication and profile management
- **API Gateway**: Centralized entry point with routing, authentication, and rate limiting
- **Swagger Documentation**: Interactive API documentation for all services
- **Monitoring and Observability**: Prometheus, Grafana, and Jaeger integration
- **Containerization**: Docker-based deployment for consistent environments
- **Kubernetes Ready**: Deployment configurations for Kubernetes

## 🏗️ Architecture

SparkFund follows a microservices architecture with the following components:

<p align="center">
  <img src="https://via.placeholder.com/800x400?text=SparkFund+Architecture" alt="SparkFund Architecture" width="800">
</p>

1. **API Gateway**: Routes requests to appropriate services, handles authentication, and rate limiting
2. **KYC Service**: Handles customer verification and onboarding with AI-powered document verification
3. **Investment Service**: Manages investments and provides AI-powered recommendations
4. **User Service**: Manages user accounts, authentication, and profiles
5. **AI Service**: Provides AI capabilities for document verification, facial recognition, and investment analysis
6. **Supporting Infrastructure**: PostgreSQL, Redis, Prometheus, Grafana, and Jaeger

## 🚀 Services

| Service            | Description                                                                                                |
| ------------------ | ---------------------------------------------------------------------------------------------------------- |
| API Gateway        | Routes requests to appropriate services, handles authentication, rate limiting, and request validation     |
| KYC Service        | Handles customer verification and onboarding with AI-powered document verification and identity validation |
| Investment Service | Manages investments, portfolios, and provides AI-powered recommendations and market analysis               |
| User Service       | Manages user accounts, authentication, profiles, and permissions                                           |
| AI Service         | Provides AI capabilities for document verification, facial recognition, NLP, and investment analysis       |

## 🚦 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/sparkfund.git
cd sparkfund
```

### Starting All Services

```bash
# Start all services
docker-compose up -d
```

This will start all services and supporting infrastructure.

### Service URLs

Once all services are running, you can access them at:

| Service            | URL                    | Description                                   |
| ------------------ | ---------------------- | --------------------------------------------- |
| API Gateway        | http://localhost:8080  | Main entry point for all services             |
| KYC Service        | http://localhost:8081  | Know Your Customer verification               |
| Investment Service | http://localhost:8082  | Investment management and AI recommendations  |
| User Service       | http://localhost:8083  | User management and authentication            |
| AI Service         | http://localhost:8001  | AI-powered document verification and analysis |
| Prometheus         | http://localhost:9090  | Metrics collection                            |
| Grafana            | http://localhost:3000  | Metrics visualization (admin/admin)           |
| Jaeger             | http://localhost:16686 | Distributed tracing                           |

### Swagger UI

All services have interactive Swagger UI documentation available:

- KYC Service: http://localhost:8081/swagger-ui.html
- Investment Service: http://localhost:8082/swagger-ui.html
- User Service: http://localhost:8083/swagger-ui.html
- AI Service: http://localhost:8001/docs

### Test Token

For local development and testing, you can use the following JWT token which has access to all endpoints:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxOTE2MjM5MDIyLCJyb2xlcyI6WyJhZG1pbiIsInVzZXIiXX0.Ks0I-dCdjWUxJEwuGP0qlyYJGXXjUYlLCRwPIZXI5Ss
```

To use this token, add it to your API requests with the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxOTE2MjM5MDIyLCJyb2xlcyI6WyJhZG1pbiIsInVzZXIiXX0.Ks0I-dCdjWUxJEwuGP0qlyYJGXXjUYlLCRwPIZXI5Ss
```

## 📚 API Documentation

All services have comprehensive Swagger documentation. You can access the interactive Swagger UI for each service at the URLs listed above.

### Key Endpoints

#### KYC Service

- `GET /health` - Health check endpoint
- `POST /api/v1/auth/login` - Login endpoint
- `POST /api/v1/verifications` - Create a verification
- `GET /api/v1/verifications` - List all verifications
- `GET /api/v1/verifications/{id}` - Get verification by ID
- `PUT /api/v1/verifications/{id}` - Update verification
- `DELETE /api/v1/verifications/{id}` - Delete verification

#### Investment Service

- `POST /api/v1/investments` - Create a new investment
- `GET /api/v1/investments/{id}` - Get investment by ID
- `GET /api/v1/investments` - List all investments
- `PUT /api/v1/investments/{id}` - Update investment
- `DELETE /api/v1/investments/{id}` - Delete investment
- `POST /api/v1/portfolios` - Create a new portfolio
- `GET /api/v1/portfolios/{id}` - Get portfolio by ID
- `PUT /api/v1/portfolios/{id}` - Update portfolio
- `DELETE /api/v1/portfolios/{id}` - Delete portfolio
- `POST /api/v1/transactions` - Create a new transaction

#### User Service

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `GET /api/v1/users` - List all users
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register

#### AI Service

- `POST /api/v1/document/analyze` - Analyze and extract information from documents
- `POST /api/v1/image/classify` - Classify images using computer vision
- `POST /api/v1/text/analyze` - Analyze text for sentiment and entities

## 🧠 AI Capabilities

SparkFund integrates advanced AI capabilities across its services:

### Investment Service AI

- **Recommendation System**: Personalized investment recommendations based on user profile and market conditions
- **Natural Language Processing**: Analysis of market news for investment signals
- **Time Series Forecasting**: Price prediction and market movement analysis
- **Reinforcement Learning**: Portfolio optimization and investment strategy learning
- **Anomaly Detection**: Fraud detection and security monitoring

### KYC Service AI

- **Document Verification**: AI-powered verification of identity documents
- **Facial Recognition**: Face matching and liveness detection
- **Risk Analysis**: AI-powered risk assessment for KYC verification

### AI Service

- **Document Analysis**: Extract information from documents using computer vision and NLP
- **Image Classification**: Classify images for various purposes
- **Text Analysis**: Sentiment analysis, entity extraction, and text classification

## 🧪 Testing

SparkFund includes comprehensive tests for all services:

### Running Tests

```bash
# Run all tests
go test ./tests/...

# Run tests for a specific service
go test ./services/kyc-service/...
go test ./services/investment-service/...
go test ./services/user-service/...
go test ./services/api-gateway/...
```

### Test Structure

Tests are organized by service and type:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test interactions between components
- **End-to-End Tests**: Test complete workflows

## 📊 Monitoring and Observability

SparkFund includes comprehensive monitoring and observability:

### Prometheus

Prometheus collects metrics from all services. Access the Prometheus UI at http://localhost:9090.

### Grafana

Grafana provides visualization of metrics with pre-configured dashboards. Access Grafana at http://localhost:3000 (username: admin, password: admin).

### Jaeger

Jaeger provides distributed tracing for request flows. Access the Jaeger UI at http://localhost:16686.

### Logging

All services use structured logging with consistent formats. Logs can be viewed using:

```bash
# View logs for all services
docker-compose logs

# View logs for a specific service
docker-compose logs [service-name]
```

## 🚢 Deployment

### Local Deployment

Local deployment is handled through Docker Compose:

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down
```

### Kubernetes Deployment

Kubernetes deployment configurations are available in the `deploy/k8s` directory.

```bash
# Apply Kubernetes configurations
kubectl apply -f deploy/k8s/
```

### Production Deployment

For production deployment, we recommend using Kubernetes with Helm charts and GitOps workflows. Configurations are available in the `deploy` directory.

## 🔧 Troubleshooting

If you encounter issues:

1. Check if all containers are running: `docker-compose ps`
2. View logs for a specific service: `docker-compose logs [service-name]`
3. Restart a specific service: `docker-compose restart [service-name]`
4. Restart all services: `docker-compose down && docker-compose up -d`

### Common Issues

- **Service Unavailable**: Check if the service container is running and healthy
- **Database Connection Issues**: Check if PostgreSQL is running and accessible
- **Authentication Failures**: Verify the JWT token is valid and not expired
- **API Gateway Routing Issues**: Check the Nginx configuration in the API Gateway service

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

<p align="center">
  Made with ❤️ by the SparkFund Team
</p>
