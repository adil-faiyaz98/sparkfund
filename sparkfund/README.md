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

# Email Service

A microservice responsible for sending emails and managing email templates in the SparkFund platform.

## Features

- Asynchronous email processing using Kafka
- Persistent storage with PostgreSQL
- Distributed tracing with Jaeger
- Metrics collection with Prometheus
- Dashboards with Grafana
- Rate limiting and request validation
- Comprehensive error handling
- Health checks and monitoring

## Prerequisites

- Docker and Docker Compose
- Go 1.19 or later
- Make (optional, for using Makefile commands)

## Environment Variables

The following environment variables can be configured:

### Server Configuration

- `PORT` - Server port (default: 8080)
- `SHUTDOWN_TIMEOUT` - Graceful shutdown timeout (default: 30s)

### Database Configuration

- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: email_service)
- `DB_SSL_MODE` - SSL mode (default: disable)

### Kafka Configuration

- `KAFKA_BROKERS` - Comma-separated list of Kafka brokers (default: localhost:9092)
- `KAFKA_TOPIC` - Kafka topic for email requests (default: email_requests)

### SMTP Configuration

- `SMTP_HOST` - SMTP server host (default: localhost)
- `SMTP_PORT` - SMTP server port (default: 587)
- `SMTP_USERNAME` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `SMTP_FROM` - Default sender email (default: noreply@example.com)

### Jaeger Configuration

- `JAEGER_ENDPOINT` - Jaeger collector endpoint (default: http://localhost:14268/api/traces)
- `JAEGER_SERVICE` - Service name for tracing (default: email-service)

### Rate Limiting

- `RATE_LIMIT` - Requests per second (default: 100)
- `RATE_LIMIT_BURST` - Burst size (default: 200)

## Development Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/adil-faiyaz98/sparkfund.git
   cd sparkfund/email-service
   ```

2. Start the development environment:

   ```bash
   make dev
   ```

3. Run database migrations:

   ```bash
   make migrate-up
   ```

4. Build and run the service:
   ```bash
   make run
   ```

## API Endpoints

### Email Operations

- `POST /api/v1/emails` - Send an email
- `GET /api/v1/emails/:id` - Get email status
- `GET /api/v1/emails` - List all emails
- `GET /api/v1/emails/stats` - Get email statistics

### Template Operations

- `POST /api/v1/templates` - Create a template
- `GET /api/v1/templates/:id` - Get a template
- `PUT /api/v1/templates/:id` - Update a template
- `DELETE /api/v1/templates/:id` - Delete a template
- `GET /api/v1/templates` - List all templates

## Monitoring

### Prometheus Metrics

The service exposes the following metrics:

- `email_sent_total` - Total number of emails sent
- `email_failed_total` - Total number of failed emails
- `email_pending_total` - Total number of pending emails
- `email_processing_duration_seconds` - Email processing duration
- `http_request_duration_seconds` - HTTP request duration
- `http_requests_total` - Total number of HTTP requests

### Distributed Tracing

The service is integrated with Jaeger for distributed tracing. You can access the Jaeger UI at http://localhost:16686.

### Grafana Dashboards

Grafana is available at http://localhost:3000 with the following default credentials:

- Username: admin
- Password: admin

## Testing

Run tests using:

```bash
make test
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
