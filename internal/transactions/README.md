# Money Pulse Transactions Service

A microservice for managing financial transactions in the Money Pulse application.

## Features

- Transaction creation and management
- Multiple transaction types (deposit, withdrawal, transfer, payment, interest, fee)
- Transaction status tracking
- Transaction categorization and tagging
- Transaction history
- Account-specific transaction views
- User-specific transaction views
- Currency support
- Metadata support for additional transaction information

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

#### Create Transaction
```http
POST /transactions
```

Request body:
```json
{
    "user_id": "uuid",
    "account_id": "uuid",
    "type": "DEPOSIT|WITHDRAWAL|TRANSFER|PAYMENT|INTEREST|FEE",
    "amount": 100.50,
    "currency": "USD",
    "description": "Monthly salary deposit",
    "category": "income",
    "tags": ["salary", "monthly"],
    "metadata": {
        "employer": "Example Corp",
        "payment_reference": "SAL-2024-03"
    },
    "source_account": "uuid",
    "destination_account": "uuid"
}
```

#### Get Transaction
```http
GET /transactions/{id}
```

#### Get User Transactions
```http
GET /transactions/user/{userId}
```

#### Get Account Transactions
```http
GET /transactions/account/{accountId}
```

#### Update Transaction Status
```http
PUT /transactions/{id}/status
```

Request body:
```json
{
    "status": "PENDING|COMPLETED|FAILED|CANCELLED",
    "error": "Error message if failed"
}
```

#### Delete Transaction
```http
DELETE /transactions/{id}
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
   export DATABASE_URL="postgresql://user:password@localhost:5432/transactions?sslmode=disable"
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
helm install transactions ./helm/transactions
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 