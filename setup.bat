@echo off

:: Create necessary directories
mkdir services\api-gateway\tmp 2>nul
mkdir services\investment-service\tmp 2>nul

:: Download dependencies for API Gateway
cd services\api-gateway
go mod tidy
cd ..\..

:: Download dependencies for Investment Service
cd services\investment-service
go mod tidy
cd ..\..

:: Create necessary directories for data
mkdir data\postgres 2>nul
mkdir data\redis 2>nul
mkdir data\grafana 2>nul

echo Setup completed successfully! 