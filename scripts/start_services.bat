@echo off
echo Starting SparkFund services...

:: Clean up any existing containers and networks
echo Cleaning up existing Docker resources...
docker-compose down --remove-orphans
docker network prune -f

:: Start Docker services in detached mode
echo Starting Docker services...
docker-compose up --build -d

:: Wait for services to start
echo Waiting for services to initialize...
timeout /t 20 /nobreak

:: Check if services are running
echo Checking service status...
docker-compose ps

echo.
echo ==========================================================
echo   SPARKFUND SERVICES INFORMATION
echo ==========================================================
echo.
echo API GATEWAY: http://localhost:8080
echo INVESTMENT SERVICE: http://localhost:8081
echo.
echo ENDPOINTS:
echo - API Info: http://localhost:8080/api
echo - Health: http://localhost:8080/health
echo - Metrics: http://localhost:8080/metrics
echo.
echo MONITORING DASHBOARDS:
echo.
echo PROMETHEUS:
echo URL: http://localhost:9090
echo Username: admin
echo Password: admin
echo.
echo GRAFANA:
echo URL: http://localhost:3000
echo Username: admin
echo Password: admin
echo.
echo ==========================================================
echo.
echo Services are starting up. You can view logs using:
echo docker-compose logs -f
echo.
echo Press Ctrl+C to stop viewing logs, services will continue running.
echo To stop all services, run: docker-compose down
echo ==========================================================

:: Show the logs
docker-compose logs -f