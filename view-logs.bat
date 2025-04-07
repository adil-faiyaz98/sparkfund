@echo off
echo ===================================================
echo SparkFund - Viewing Service Logs
echo ===================================================

if "%1"=="" (
    echo Viewing logs for all services...
    docker-compose logs -f
) else (
    echo Viewing logs for %1...
    docker-compose logs -f %1
)

echo ===================================================
