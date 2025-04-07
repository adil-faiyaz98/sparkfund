@echo off
echo ===================================================
echo SparkFund - Starting All Microservices
echo ===================================================

echo Creating monitoring directories if they don't exist...
if not exist "monitoring\prometheus" mkdir monitoring\prometheus
if not exist "monitoring\grafana\provisioning\datasources" mkdir monitoring\grafana\provisioning\datasources
if not exist "monitoring\grafana\provisioning\dashboards" mkdir monitoring\grafana\provisioning\dashboards
if not exist "monitoring\grafana\dashboards" mkdir monitoring\grafana\dashboards

echo Creating Prometheus configuration...
(
echo global:
echo   scrape_interval: 15s
echo   evaluation_interval: 15s
echo scrape_configs:
echo   - job_name: 'prometheus'
echo     static_configs:
echo       - targets: ['localhost:9090']
echo   - job_name: 'api-gateway'
echo     static_configs:
echo       - targets: ['api-gateway:8080']
echo   - job_name: 'kyc-service'
echo     static_configs:
echo       - targets: ['kyc-service:8081']
echo   - job_name: 'investment-service'
echo     static_configs:
echo       - targets: ['investment-service:8082']
echo   - job_name: 'user-service'
echo     static_configs:
echo       - targets: ['user-service:8083']
echo   - job_name: 'ai-service'
echo     static_configs:
echo       - targets: ['ai-service:8001']
) > monitoring\prometheus\prometheus.yml

echo Creating Grafana datasource configuration...
(
echo apiVersion: 1
echo datasources:
echo   - name: Prometheus
echo     type: prometheus
echo     access: proxy
echo     url: http://prometheus:9090
echo     isDefault: true
) > monitoring\grafana\provisioning\datasources\datasource.yml

echo Creating Grafana dashboard configuration...
(
echo apiVersion: 1
echo providers:
echo   - name: 'default'
echo     orgId: 1
echo     folder: ''
echo     type: file
echo     disableDeletion: false
echo     updateIntervalSeconds: 10
echo     options:
echo       path: /var/lib/grafana/dashboards
) > monitoring\grafana\provisioning\dashboards\dashboard.yml

echo Building and starting all services...
docker-compose -f docker-compose-simple.yml up --build -d

echo Waiting for services to start...
timeout /t 30 /nobreak > nul

echo Checking service status...
docker-compose -f docker-compose-all.yml ps

echo ===================================================
echo SparkFund Services:
echo ===================================================
echo API Gateway:         http://localhost:8080
echo KYC Service:         http://localhost:8081
echo Investment Service:  http://localhost:8082
echo User Service:        http://localhost:8083
echo AI Service:          http://localhost:8001
echo ===================================================
echo Monitoring:
echo ===================================================
echo Prometheus:          http://localhost:9090
echo Grafana:             http://localhost:3000 (admin/admin)
echo Jaeger:              http://localhost:16686
echo ===================================================

echo To view logs, run: docker-compose -f docker-compose-all.yml logs -f
echo To stop all services, run: docker-compose -f docker-compose-all.yml down
