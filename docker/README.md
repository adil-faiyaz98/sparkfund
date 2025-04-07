# Docker Compose Configuration for SparkFund

This directory contains Docker Compose configurations for the SparkFund application.

## Docker Compose Structure

The SparkFund application uses a hierarchical Docker Compose structure:

1. **Root `docker-compose.yml`**: Runs all services together
2. **Service-specific Docker Compose files**: Run individual services independently

### Root Docker Compose

The root `docker-compose.yml` file in the project root directory runs all services together:

- API Gateway
- Auth Service
- KYC Service
- Investment Service
- AML Service
- Security Monitoring
- Postgres
- Redis
- Elasticsearch

### Service-Specific Docker Compose

Each service has its own Docker Compose files:

- **Production**: `docker-compose.yml`
- **Development**: `docker-compose.dev.yml`

## Usage

### Running All Services

```bash
# From the project root
docker-compose up -d
```

### Running Individual Services

#### KYC Service

```bash
# Production
cd services/kyc-service
docker-compose up -d

# Development
cd services/kyc-service
docker-compose -f docker-compose.dev.yml up -d
```

#### Investment Service

```bash
# Production
cd services/investment-service
docker-compose up -d

# Development
cd services/investment-service
docker-compose -f docker-compose.dev.yml up -d
```

## Port Mapping

To avoid port conflicts, each service uses different ports:

| Service | Production Port | Development Port |
|---------|----------------|------------------|
| API Gateway | 8080 | 8080 |
| KYC Service | 8081 | 8082 |
| Investment Service | 8082 | 8083 |
| Auth Service | 8084 | 8085 |
| AML Service | 8086 | 8087 |

## Volumes

The Docker Compose files define several volumes:

- **`postgres-data`**: Stores PostgreSQL data
- **`redis-data`**: Stores Redis data
- **`kyc-data`**: Stores KYC documents
- **`investment-data`**: Stores investment data
- **`elastic-data`**: Stores Elasticsearch data

## Networks

The services use the default network created by Docker Compose.

## Environment Variables

See the Docker Compose files for the environment variables used by each service.
