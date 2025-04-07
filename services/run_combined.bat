@echo off
echo Building and running Docker Compose setup...
docker-compose -f docker-compose-combined.yml up --build -d

echo Waiting for services to start...
timeout /t 20 /nobreak > nul

echo Checking if services are running...
docker-compose -f docker-compose-combined.yml ps

echo Showing logs...
docker-compose -f docker-compose-combined.yml logs

echo KYC service is running at http://localhost:8081
echo KYC API documentation is available at http://localhost:8081/swagger-ui/index.html
echo KYC Health endpoint is available at http://localhost:8081/health
echo.
echo AI service is running at http://localhost:8001
echo AI API documentation is available at http://localhost:8001/docs
echo AI Health endpoint is available at http://localhost:8001/health
