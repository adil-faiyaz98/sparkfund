@echo off
echo Starting all SparkFund services...

cd %~dp0..
docker-compose up -d

echo All services started successfully!
echo API Gateway: http://localhost:8080
echo KYC Service: http://localhost:8081
echo Investment Service: http://localhost:8082
echo Auth Service: http://localhost:8084
echo AML Service: http://localhost:8086
