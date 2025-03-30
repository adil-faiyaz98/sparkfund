# AML Service

The AML (Anti-Money Laundering) Service is a microservice responsible for monitoring and analyzing financial transactions to detect potential money laundering activities. It provides risk assessment, transaction screening, and alert generation capabilities.

## Features

- Transaction processing and risk assessment
- Sanctions list screening
- PEP (Politically Exposed Person) screening
- Watch list monitoring
- Risk profile management
- Alert generation and management
- Transaction history tracking

## Prerequisites

- Go 1.21 or later
- PostgreSQL 13 or later
- Docker (optional, for containerized deployment)

## Configuration

The service can be configured using a YAML file (`config.yaml`) or environment variables. The following configuration options are available:

```yaml
server:
  port: 8080
  readTimeout: 10
  writeTimeout: 10
  idleTimeout: 120

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: aml_service
  sslmode: disable

security:
  cors:
    allowedOrigins:
      - "*"
    allowedMethods:
      - GET
      - POST
      - PUT
      - DELETE
      - OPTIONS
    allowedHeaders:
      - Content-Type
      - Authorization
    exposedHeaders:
      - Content-Length
    allowCredentials: true
    maxAge: 3600
```

## API Endpoints

### Transactions

- `POST /api/v1/aml/transactions` - Process a new transaction
- `GET /api/v1/aml/transactions` - List transactions with optional filters
- `GET /api/v1/aml/transactions/{id}` - Get transaction details
- `POST /api/v1/aml/transactions/{id}/flag` - Flag a transaction for review
- `POST /api/v1/aml/transactions/{id}/review` - Review a flagged transaction

### Risk Profiles

- `GET /api/v1/aml/risk-profiles/{userId}` - Get user's risk profile

## Development

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up the database:
   ```bash
   createdb aml_service
   ```
4. Run the service:
   ```bash
   go run cmd/main.go
   ```

## Testing

Run the tests:
```bash
go test ./...
```

## Docker

Build the Docker image:
```bash
docker build -t aml-service .
```

Run the container:
```bash
docker run -p 8080:8080 aml-service
```

## Metrics

The service exposes Prometheus metrics at `/metrics` endpoint. Key metrics include:

- `aml_transaction_process_time_seconds` - Time taken to process transactions
- `aml_errors_total` - Total number of errors by type

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 