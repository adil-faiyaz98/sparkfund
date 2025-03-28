# KYC Service

A microservice for handling Know Your Customer (KYC) verification processes in the SparkFund platform.

## Features

- Submit KYC documents
- Check KYC status
- Verify KYC submissions
- Reject KYC submissions with reasons
- List pending KYC submissions
- Health check endpoint

## Prerequisites

- Go 1.21 or later
- PostgreSQL 13 or later
- Docker (optional)

## Configuration

The service can be configured using environment variables or a `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sparkfund
SERVER_PORT=8080
```

## API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
    "status": "ok"
}
```

### Submit KYC
```
POST /api/v1/kyc
```

Request body:
```json
{
    "firstName": "John",
    "lastName": "Doe",
    "dateOfBirth": "1990-01-01",
    "address": "123 Main St",
    "city": "New York",
    "country": "US",
    "postalCode": "10001",
    "documentType": "passport",
    "documentNumber": "123456789",
    "documentFront": "base64_encoded_image",
    "documentBack": "base64_encoded_image",
    "selfieImage": "base64_encoded_image"
}
```

### Get KYC Status
```
GET /api/v1/kyc/{id}
```

### Verify KYC
```
POST /api/v1/kyc/{id}/verify
```

### Reject KYC
```
POST /api/v1/kyc/{id}/reject
```

Request body:
```json
{
    "reason": "Document quality is poor"
}
```

### List Pending KYC
```
GET /api/v1/kyc/pending
```

## Development

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
3. Run `go mod download` to install dependencies
4. Run `go run cmd/server/main.go` to start the service

## Docker

Build the image:
```bash
docker build -t sparkfund/kyc-service .
```

Run the container:
```bash
docker run -p 8080:8080 --env-file .env sparkfund/kyc-service
```

## Testing

Run tests:
```bash
go test ./...
```

## License

MIT 