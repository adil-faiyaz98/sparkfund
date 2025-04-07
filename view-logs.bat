@echo off
echo ===================================================
echo SparkFund - Viewing Service Logs
echo ===================================================

if "%1"=="" (
    echo Viewing logs for all services...
    docker-compose -f docker-compose-simple.yml logs -f
) else (
    echo Viewing logs for %1...
    docker-compose -f docker-compose-simple.yml logs -f %1
)

echo ===================================================
