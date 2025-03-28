# SparkFund Microservices

This repository contains the microservices for the SparkFund platform, including AML (Anti-Money Laundering) and KYC (Know Your Customer) services.

## Services

### AML Service
The AML service handles anti-money laundering checks, alerts, and watchlist management.

### KYC Service
The KYC service manages customer verification and due diligence processes.

## Development Setup

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or later
- Make
- Helm 3.x
- kubectl

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/sparkfund/sparkfund.git
cd sparkfund
```

2. Start the development environment:
```bash
make dev
```

This will start:
- AML Service (port 8080)
- KYC Service (port 8081)
- PostgreSQL (port 5432)
- Prometheus (port 9090)
- Grafana (port 3000)
- Jaeger (port 16686)

### Accessing Services

- AML Service: http://localhost:8080
- KYC Service: http://localhost:8081
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- Jaeger UI: http://localhost:16686

## Deployment

### Kubernetes Deployment

1. Create a namespace:
```bash
kubectl create namespace sparkfund
```

2. Deploy the services:
```bash
make deploy
```

This will deploy:
- AML Service
- KYC Service
- Prometheus
- Grafana
- Jaeger

### Monitoring

The services are configured with:
- Prometheus metrics
- Grafana dashboards
- Jaeger tracing
- Health checks
- Resource limits and requests

## Development Commands

```bash
# Build all services
make build

# Run tests
make test

# Deploy to Kubernetes
make deploy

# Clean up
make clean

# Generate protobuf
make proto

# Run locally
make run

# Run with hot reload
make dev

# Check code style
make lint

# Generate swagger docs
make swagger

# Run database migrations
make migrate

# Run security scan
make security
```

## Project Structure

```
sparkfund/
├── aml-service/           # AML Service
├── kyc-service/          # KYC Service
├── helm/                 # Helm charts
├── prometheus/           # Prometheus configuration
├── grafana/             # Grafana dashboards
├── jaeger/              # Jaeger configuration
├── docker-compose.yml   # Production compose file
├── docker-compose.dev.yml # Development compose file
└── Makefile             # Build and deployment commands
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 