# SparkFund Investment Platform

A microservices-based investment platform built with Go, featuring an API Gateway and Investment Service.

## Project Structure

```
.
├── services/
│   ├── api-gateway/        # API Gateway service
│   └── investment-service/ # Investment Service
├── deploy/
│   └── k8s/               # Kubernetes deployment configurations
└── docker-compose.yml     # Local development setup
```

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL
- Redis
- Make (optional, for using Makefile commands)

## Local Development

1. Clone the repository:
```bash
git clone https://github.com/adil-faiyaz98/sparkfund.git
cd sparkfund
```

2. Start the local development environment:
```bash
docker-compose up --build
```

This will start all required services:
- API Gateway (http://localhost:8080)
- Investment Service (http://localhost:8081)
- PostgreSQL (localhost:5432)
- Redis (localhost:6379)
- Prometheus (http://localhost:9090)
- Grafana (http://localhost:3000)

## API Endpoints

### Investment Service

- `POST /api/v1/investments` - Create a new investment
- `GET /api/v1/investments/:id` - Get investment by ID
- `GET /api/v1/investments` - List all investments
- `PUT /api/v1/investments/:id` - Update investment
- `DELETE /api/v1/investments/:id` - Delete investment

### Portfolio Management

- `POST /api/v1/portfolios` - Create a new portfolio
- `GET /api/v1/portfolios/:id` - Get portfolio by ID
- `PUT /api/v1/portfolios/:id` - Update portfolio
- `DELETE /api/v1/portfolios/:id` - Delete portfolio

### Transactions

- `POST /api/v1/transactions` - Create a new transaction

## Testing

### Running Tests

```bash
# Run tests for all services
make test

# Run tests for specific service
cd services/api-gateway && go test ./...
cd services/investment-service && go test ./...
```

### API Testing

You can use curl or any API testing tool like Postman to test the endpoints. Here's an example:

```bash
# Create a new investment
curl -X POST http://localhost:8080/api/v1/investments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123",
    "amount": 1000,
    "type": "STOCK",
    "symbol": "AAPL",
    "quantity": 10
  }'

# Get all investments
curl http://localhost:8080/api/v1/investments
```

## Monitoring

- Prometheus metrics: http://localhost:9090
- Grafana dashboards: http://localhost:3000 (username: admin, password: admin)

## Development Guidelines

1. Follow Go best practices and coding standards
2. Write tests for new features
3. Update documentation when making changes
4. Use meaningful commit messages

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
