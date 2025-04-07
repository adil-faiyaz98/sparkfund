@echo off
echo ===================================================
echo SparkFund - Checking Service Status
echo ===================================================

echo Checking container status...
docker-compose -f docker-compose-simple.yml ps

echo ===================================================
echo Checking API Gateway health...
curl -s http://localhost:8080/health || echo API Gateway is not responding

echo ===================================================
echo Checking KYC Service health...
curl -s http://localhost:8081/health || echo KYC Service is not responding

echo ===================================================
echo Checking Investment Service health...
curl -s http://localhost:8082/health || echo Investment Service is not responding

echo ===================================================
echo Checking User Service health...
curl -s http://localhost:8083/health || echo User Service is not responding

echo ===================================================
echo Checking AI Service health...
curl -s http://localhost:8001/health || echo AI Service is not responding

echo ===================================================
