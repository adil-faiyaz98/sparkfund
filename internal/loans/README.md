# Money Pulse Loans Service

A microservice for managing loans and credit facilities in the Money Pulse application.

## Features

- Loan application processing
- Credit scoring and risk assessment
- Loan approval workflow
- Payment scheduling and tracking
- Interest rate calculation
- Loan status monitoring
- AI-powered credit risk assessment

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
All endpoints require Bearer token authentication:
```
Authorization: Bearer <your_jwt_token>
```

### Endpoints

#### Create Loan Application
```http
POST /loans
```

Request body:
```json
{
    "user_id": "uuid",
    "account_id": "uuid",
    "loan_type": "PERSONAL|MORTGAGE|BUSINESS|EDUCATION",
    "amount": "float",
    "term_months": "integer",
    "purpose": "string",
    "interest_rate": "float"
}
```

#### Get Loan
```http
GET /loans/{id}
```

#### Get User Loans
```http
GET /loans/user/{userId}
```

#### Get Account Loans
```http
GET /loans/account/{accountId}
```

#### Update Loan Status
```http
PUT /loans/{id}/status
```

Request body:
```json
{
    "status": "PENDING|APPROVED|REJECTED|ACTIVE|PAID|DEFAULTED",
    "notes": "string"
}
```

#### Make Payment
```http
POST /loans/{id}/payments
```

Request body:
```json
{
    "amount": "float",
    "payment_date": "timestamp"
}
```

#### Get Loan Payments
```http
GET /loans/{id}/payments
```

## Development

### Prerequisites
- Go 1.24 or higher
- PostgreSQL 13 or higher
- Docker (optional)

### Setup
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables:
   ```bash
   export DATABASE_URL="postgresql://user:password@localhost:5432/loans?sslmode=disable"
   ```
4. Run migrations:
   ```bash
   go run cmd/migrate/main.go
   ```
5. Start the service:
   ```bash
   go run cmd/main.go
   ```

### Running Tests
```bash
go test ./...
```

## Deployment

The service can be deployed using the provided Helm chart:

```bash
helm install loans ./helm/loans
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 