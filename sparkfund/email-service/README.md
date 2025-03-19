# Email Service

A microservice responsible for sending emails and managing email templates. This service provides a RESTful API for sending emails, managing email templates, and tracking email delivery status.

## Features

- Send emails with support for CC, BCC, and attachments
- Create and manage email templates with variable substitution
- Track email delivery status and history
- RESTful API with Swagger documentation
- PostgreSQL database for persistent storage
- Kafka integration for asynchronous email processing
- Prometheus metrics and Grafana dashboards
- Jaeger tracing for distributed tracing

## Prerequisites

- Go 1.19 or later
- Docker and Docker Compose
- PostgreSQL 14 or later
- Kafka (provided via Docker Compose)
- Jaeger (provided via Docker Compose)
- Prometheus (provided via Docker Compose)
- Grafana (provided via Docker Compose)

## Environment Variables

The following environment variables can be configured:

```env
# Server Configuration
PORT=8080
SHUTDOWN_TIMEOUT=10s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=email_service
DB_SSL_MODE=disable

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=emails

# SMTP Configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
SMTP_FROM=noreply@example.com

# Jaeger Configuration
JAEGER_ENDPOINT=http://localhost:14268/api/traces
JAEGER_SERVICE=email-service

# Rate Limiting
RATE_LIMIT=100
RATE_LIMIT_BURST=10
```

## Development Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/adil-faiyaz98/sparkfund/email-service.git
   cd email-service
   ```

2. Start the development environment:

   ```bash
   docker-compose up -d
   ```

3. Run database migrations:

   ```bash
   make migrate-up
   ```

4. Build the service:

   ```bash
   make build
   ```

5. Run the service:
   ```bash
   make run
   ```

## API Endpoints

### Email Operations

- `POST /api/v1/emails` - Send an email
- `GET /api/v1/emails` - Get email logs

### Template Operations

- `POST /api/v1/templates` - Create a new template
- `GET /api/v1/templates/:id` - Get a template
- `PUT /api/v1/templates/:id` - Update a template
- `DELETE /api/v1/templates/:id` - Delete a template

### Health Check

- `GET /health` - Health check endpoint

### API Documentation

- `GET /swagger/*` - Swagger UI documentation

## Monitoring

### Prometheus Metrics

The service exposes Prometheus metrics at `/metrics`. Key metrics include:

- `email_sent_total` - Total number of emails sent
- `email_failed_total` - Total number of failed email attempts
- `template_operations_total` - Total number of template operations
- `http_requests_total` - Total number of HTTP requests
- `http_request_duration_seconds` - HTTP request duration

### Jaeger Tracing

Distributed tracing is available through Jaeger. Access the Jaeger UI at `http://localhost:16686`.

### Grafana Dashboards

Pre-built Grafana dashboards are available at `http://localhost:3000`.

## Testing

Run tests:

```bash
make test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
