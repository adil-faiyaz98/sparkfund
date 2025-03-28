#!/bin/bash

# Create necessary directories
mkdir -p services/api-gateway/tmp
mkdir -p services/investment-service/tmp

# Download dependencies for API Gateway
cd services/api-gateway
go mod tidy
cd ../..

# Download dependencies for Investment Service
cd services/investment-service
go mod tidy
cd ../..

# Create necessary directories for data
mkdir -p data/postgres
mkdir -p data/redis
mkdir -p data/grafana

# Set proper permissions
chmod -R 777 data

echo "Setup completed successfully!" 