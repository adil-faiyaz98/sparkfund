@echo off
echo Starting Investment Service in development mode...

cd %~dp0..\services\investment-service
docker-compose -f docker-compose.dev.yml up -d

echo Investment Service started successfully!
echo Investment Service: http://localhost:8083
