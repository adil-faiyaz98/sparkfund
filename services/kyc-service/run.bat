@echo off
echo Building and running Docker Compose setup...
docker-compose up --build -d

echo Waiting for services to start...
timeout /t 10 /nobreak > nul

echo Checking if services are running...
docker-compose ps

echo Showing logs...
docker-compose logs

echo KYC service is running at http://localhost:8080
echo API documentation is available at http://localhost:8080/swagger-ui/index.html
echo Health endpoint is available at http://localhost:8080/health
