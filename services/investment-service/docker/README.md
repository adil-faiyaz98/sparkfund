# Docker Compose Configuration for Investment Service

This directory contains Docker Compose configurations for the Investment service.

## Docker Compose Files

- **`docker-compose.yml`**: Production configuration for running the Investment service independently
- **`docker-compose.dev.yml`**: Development configuration for running the Investment service independently

## Usage

### Running in Production Mode

```bash
cd services/investment-service
docker-compose up -d
```

### Running in Development Mode

```bash
cd services/investment-service
docker-compose -f docker-compose.dev.yml up -d
```

## Integration with Root Docker Compose

The Investment service is also included in the root `docker-compose.yml` file, which runs all services together. To run the entire application:

```bash
# From the project root
docker-compose up -d
```

## Volumes

- **`postgres_data`**: Stores PostgreSQL data
- **`investment-data`**: Stores investment data
- **`redis-data`**: Stores Redis data (development only)

## Networks

The services use the default network created by Docker Compose.

## Environment Variables

See the Docker Compose files for the environment variables used by each service.
