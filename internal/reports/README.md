# Money Pulse Reports Service

A microservice for generating and managing financial reports in the Money Pulse application.

## Features

- Balance sheet reports
- Transaction history reports
- Investment portfolio reports
- Loan payment reports
- Tax reports
- Multiple report formats (PDF, CSV, JSON, XLSX)
- Asynchronous report generation
- Report status tracking
- AI-powered report insights

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

#### Create Report Request
```http
POST /reports
```

Request body:
```json
{
    "user_id": "uuid",
    "type": "BALANCE|TRANSACTION|INVESTMENT|LOAN|TAX",
    "format": "PDF|CSV|JSON|XLSX",
    "parameters": {
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "account_id": "uuid",
        "categories": ["income", "expenses"]
    }
}
```

#### Get Report
```http
GET /reports/{id}
```

#### Get User Reports
```http
GET /reports/user/{userId}
```

#### Update Report Status
```http
PUT /reports/{id}/status
```

Request body:
```json
{
    "status": "PENDING|GENERATING|COMPLETED|FAILED",
    "file_url": "https://storage.example.com/reports/report.pdf",
    "error": "Error message if failed"
}
```

#### Delete Report
```http
DELETE /reports/{id}
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
   export DATABASE_URL="postgresql://user:password@localhost:5432/reports?sslmode=disable"
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
helm install reports ./helm/reports
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 