# B2B Banking Microservices Platform

A comprehensive B2B microservices platform for banking services including KYC, AML, fraud detection, and more.

## Architecture

This project follows a microservices architecture using Go. Each service is independently deployable and has its own:

- API endpoints
- Database access layer
- Business logic
- Configuration

## Services

- **API Gateway**: Central entry point for all API requests
- **KYC Service**: Handles Know Your Customer verification processes
- **AML Service**: Anti-Money Laundering compliance and monitoring
- **Fraud Detection**: Real-time and batch fraud detection
- **Credit Scoring**: Credit assessment and scoring
- **Risk Management**: Financial forecasting and risk modeling
- **Notification**: Customer communication (email, SMS, etc.)
- **Consent Management**: Handling customer consent for data usage
- **Logging**: Centralized logging service
- **Security**: Authentication, encryption, and security features
- **Email**: Email handling and delivery

## Development

### Prerequisites

- Go 1.20+
- Docker and Docker Compose
- Make

### Getting Started

1. Clone the repository
2. Run `make docker-up` to start all services
3. Access the API Gateway at http://localhost:8000

### Common Commands

- `make build`: Build all services
- `make test`: Run tests for all services
- `make docker-up`: Start all services with Docker
- `make docker-down`: Stop all services
- `make lint`: Run linters on all services
- `make new-service`: Create a new service from template
