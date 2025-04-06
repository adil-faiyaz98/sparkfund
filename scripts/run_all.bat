@echo off
echo Building and running Docker Compose setup...
docker-compose -f services/docker-compose.yml up --build -d

echo Waiting for services to start...
timeout /t 20 /nobreak > nul

echo Checking if services are running...
docker-compose -f services/docker-compose.yml ps

echo Showing logs...
docker-compose -f services/docker-compose.yml logs

echo KYC service is running at http://localhost:8080
echo KYC API documentation is available at http://localhost:8080/swagger-ui/index.html
echo KYC Health endpoint is available at http://localhost:8080/health
echo.
echo AI service is running at http://localhost:8000
echo AI API documentation is available at http://localhost:8000/docs
echo AI Health endpoint is available at http://localhost:8000/health
