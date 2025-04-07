# SparkFund Local Deployment

This guide provides instructions for deploying and running all SparkFund microservices locally.

## Quick Start

### Windows

1. Run all services:
   ```
   .\run-all-services.bat
   ```

2. Check service status:
   ```
   .\check-services.bat
   ```

3. View logs:
   ```
   .\view-logs.bat
   ```

4. Stop all services:
   ```
   .\stop-all-services.bat
   ```

### Linux/Mac

1. Make scripts executable:
   ```
   chmod +x *.sh
   ```

2. Run all services:
   ```
   ./run-all-services.sh
   ```

3. Check service status:
   ```
   ./check-services.sh
   ```

4. View logs:
   ```
   ./view-logs.sh
   ```

5. Stop all services:
   ```
   ./stop-all-services.sh
   ```

## Service URLs

Once all services are running, you can access them at the following URLs:

- **API Gateway**: http://localhost:8080
- **KYC Service**: http://localhost:8081
- **Investment Service**: http://localhost:8082
- **User Service**: http://localhost:8083
- **AI Service**: http://localhost:8001

## Monitoring

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

## Detailed Documentation

For more detailed information, please refer to the [LOCAL-DEPLOYMENT.md](LOCAL-DEPLOYMENT.md) file.
