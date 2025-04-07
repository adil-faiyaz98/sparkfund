# SparkFund Local Deployment Guide

This guide provides instructions for deploying and managing all SparkFund microservices locally.

## Prerequisites

- Docker and Docker Compose installed
- PowerShell (for Windows) or Bash (for Linux/Mac)
- Git (to clone the repository)

## Services Overview

The SparkFund platform consists of the following services:

1. **API Gateway** - Entry point for all client requests (Port: 8080)
2. **KYC Service** - Know Your Customer functionality (Port: 8081)
3. **Investment Service** - Investment management (Port: 8082)
4. **User Service** - User management (Port: 8083)
5. **AI Service** - AI-powered document verification and analysis (Port: 8001)

Supporting infrastructure:
- **PostgreSQL** - Primary database (Port: 5432)
- **Redis** - Caching and session management (Port: 6379)
- **Prometheus** - Metrics collection (Port: 9090)
- **Grafana** - Metrics visualization (Port: 3000)
- **Jaeger** - Distributed tracing (Port: 16686)

## Deployment Instructions

### Starting All Services

To start all services, run the following command:

```bash
./run-all-services.bat  # Windows
# or
./run-all-services.sh   # Linux/Mac
```

This script will:
1. Create necessary directories for monitoring
2. Generate configuration files for Prometheus and Grafana
3. Build and start all services using Docker Compose
4. Display the URLs for accessing each service

### Checking Service Status

To check the status of all services, run:

```bash
./check-services.bat  # Windows
# or
./check-services.sh   # Linux/Mac
```

### Viewing Service Logs

To view logs for all services:

```bash
./view-logs.bat  # Windows
# or
./view-logs.sh   # Linux/Mac
```

To view logs for a specific service:

```bash
./view-logs.bat api-gateway  # Windows
# or
./view-logs.sh api-gateway   # Linux/Mac
```

### Testing Services

To test all services:

```bash
./test-services.ps1  # Windows
# or
./test-services.sh   # Linux/Mac
```

### Stopping All Services

To stop all services:

```bash
./stop-all-services.bat  # Windows
# or
./stop-all-services.sh   # Linux/Mac
```

## Accessing Services

Once all services are running, you can access them at the following URLs:

- **API Gateway**: http://localhost:8080
- **KYC Service**: http://localhost:8081
- **Investment Service**: http://localhost:8082
- **User Service**: http://localhost:8083
- **AI Service**: http://localhost:8001

Monitoring:
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

## Troubleshooting

### Common Issues

1. **Port conflicts**: If you have other services running on the same ports, you'll need to modify the port mappings in the `docker-compose-all.yml` file.

2. **Docker errors**: Make sure Docker is running and you have sufficient permissions.

3. **Service dependencies**: Some services depend on others. If a service is failing to start, check the logs to see if it's waiting for a dependency.

### Viewing Detailed Logs

For more detailed logs, use:

```bash
docker-compose -f docker-compose-all.yml logs -f [service_name]
```

### Restarting a Specific Service

To restart a specific service:

```bash
docker-compose -f docker-compose-all.yml restart [service_name]
```

## Development Workflow

1. Make changes to the code
2. Rebuild and restart the affected service:
   ```bash
   docker-compose -f docker-compose-all.yml up -d --build [service_name]
   ```
3. Test your changes
4. Commit your changes to Git

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Prometheus Documentation](https://prometheus.io/docs/introduction/overview/)
- [Grafana Documentation](https://grafana.com/docs/)
