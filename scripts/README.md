# SparkFund Scripts

This directory contains various helper scripts for managing the SparkFund platform.

## Web Interface

- `swagger-hub.html` - A web interface for accessing all service Swagger documentation

## Service Management Scripts

These scripts are stored in the Git repository but may not be present in your local directory after cleanup:

### Windows Scripts (.bat)

- `run-all-services.bat` - Start all services using Docker Compose
- `stop-all-services.bat` - Stop all services
- `view-logs.bat` - View logs for all services or a specific service
- `test-endpoints.ps1` - Test all API endpoints

### Unix/Linux/Mac Scripts (.sh)

- `run-all-services.sh` - Start all services using Docker Compose
- `stop-all-services.sh` - Stop all services
- `view-logs.sh` - View logs for all services or a specific service
- `test-endpoints.sh` - Test all API endpoints

### PowerShell Scripts (.ps1)

- `check-services.ps1` - Check the health of all services

## Usage

### Windows

```powershell
# Start all services
.\scripts\run-all-services.bat

# Stop all services
.\scripts\stop-all-services.bat

# View logs for all services
.\scripts\view-logs.bat

# View logs for a specific service (e.g., kyc-service)
.\scripts\view-logs.bat kyc-service

# Test all endpoints
.\scripts\test-endpoints.ps1

# Check service health
.\scripts\check-services.ps1
```

### Unix/Linux/Mac

```bash
# Start all services
./scripts/run-all-services.sh

# Stop all services
./scripts/stop-all-services.sh

# View logs for all services
./scripts/view-logs.sh

# View logs for a specific service (e.g., kyc-service)
./scripts/view-logs.sh kyc-service

# Test all endpoints
./scripts/test-endpoints.sh

# Check service health
./scripts/check-services.sh
```

## Web Interface Usage

To use the Swagger Hub web interface:

1. Open the `swagger-hub.html` file in a web browser
2. Click on the links to access the Swagger UI for each service
3. Use the provided test JWT token for authentication
