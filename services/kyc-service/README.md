# KYC Service

The KYC (Know Your Customer) Service is a microservice responsible for managing customer identity verification and document processing in the SparkFund platform. It handles document uploads, verification processes, and maintains KYC profiles for users.

## Features

- Document Management
  - Upload and store KYC documents
  - Support for multiple document types (ID, proof of address, etc.)
  - Document status tracking and expiration management
  - Secure document storage with encryption

- KYC Profile Management
  - Create and update customer profiles
  - Risk level assessment
  - Profile status tracking
  - Comprehensive customer information storage

- Verification Process
  - Document verification workflow
  - Multiple verification methods support
  - Verification history tracking
  - Confidence scoring system

- Integration
  - AML service integration for enhanced screening
  - Notification service integration for status updates
  - User service integration for profile management

## API Endpoints

### Documents
- `POST /api/v1/documents` - Upload a new document
- `GET /api/v1/documents/:id` - Get document details
- `GET /api/v1/documents` - List documents
- `PUT /api/v1/documents/:id/status` - Update document status
- `DELETE /api/v1/documents/:id` - Delete document
- `GET /api/v1/documents/pending` - List pending documents
- `GET /api/v1/documents/expired` - List expired documents

### Profiles
- `POST /api/v1/profiles` - Create a new KYC profile
- `GET /api/v1/profiles/:user_id` - Get profile details
- `PUT /api/v1/profiles/:user_id` - Update profile
- `PUT /api/v1/profiles/:user_id/status` - Update profile status
- `PUT /api/v1/profiles/:user_id/risk-level` - Update risk level
- `GET /api/v1/profiles` - List profiles
- `GET /api/v1/profiles/stats` - Get profile statistics
- `DELETE /api/v1/profiles/:user_id` - Delete profile

### Verifications
- `POST /api/v1/verifications` - Create verification record
- `GET /api/v1/verifications/:id` - Get verification details
- `GET /api/v1/verifications/document/:document_id` - List document verifications
- `PUT /api/v1/verifications/:id` - Update verification
- `DELETE /api/v1/verifications/:id` - Delete verification
- `GET /api/v1/verifications/stats` - Get verification statistics
- `GET /api/v1/verifications/document/:document_id/history` - Get verification history
- `GET /api/v1/verifications/document/:document_id/summary` - Get verification summary

## Prerequisites

- Go 1.21 or later
- PostgreSQL 14 or later
- Redis (optional, for caching)
- Docker (optional, for containerization)

## Configuration

The service can be configured using environment variables or a `.env` file. See `.env.example` for available configuration options.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/sparkfund/kyc-service.git
cd kyc-service
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the service:
```bash
go run cmd/server/main.go
```

## Docker Deployment

1. Build the Docker image:
```bash
docker build -t sparkfund/kyc-service .
```

2. Run the container:
```bash
docker run -p 8080:8080 \
  --env-file .env \
  -v /path/to/documents:/data/documents \
  sparkfund/kyc-service
```

## Development

### Running Tests
```bash
go test -v ./...
```

### Code Style
The project follows the standard Go code style. Run the linter:
```bash
golangci-lint run
```

### API Documentation
The API documentation is available in Swagger format at `/swagger/index.html` when running the service.

## Security

- All endpoints require authentication
- Documents are encrypted at rest
- File uploads are validated and sanitized
- Rate limiting is implemented
- CORS is configured for secure cross-origin requests

## Monitoring

The service exposes metrics at `/metrics` for Prometheus integration and health checks at `/health`.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 