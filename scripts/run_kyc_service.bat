@echo off
echo Starting KYC Service in development mode...

cd %~dp0..\services\kyc-service
docker-compose -f docker-compose.dev.yml up -d

echo KYC Service started successfully!
echo KYC Service: http://localhost:8082
echo KYC ML Service: http://localhost:8000
