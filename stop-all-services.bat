@echo off
echo ===================================================
echo SparkFund - Stopping All Microservices
echo ===================================================

echo Stopping and removing all containers...
docker-compose -f docker-compose-simple.yml down

echo ===================================================
echo All services have been stopped.
echo ===================================================
