#!/bin/bash
echo "==================================================="
echo "SparkFund - Starting All Microservices"
echo "==================================================="

echo "Creating monitoring directories if they don't exist..."
mkdir -p monitoring/prometheus
mkdir -p monitoring/grafana/provisioning/datasources
mkdir -p monitoring/grafana/provisioning/dashboards
mkdir -p monitoring/grafana/dashboards

echo "Creating Prometheus configuration..."
cat > monitoring/prometheus/prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'api-gateway'
    static_configs:
      - targets: ['api-gateway:8080']
  - job_name: 'kyc-service'
    static_configs:
      - targets: ['kyc-service:8081']
  - job_name: 'investment-service'
    static_configs:
      - targets: ['investment-service:8082']
  - job_name: 'user-service'
    static_configs:
      - targets: ['user-service:8083']
  - job_name: 'ai-service'
    static_configs:
      - targets: ['ai-service:8001']
EOF

echo "Creating Grafana datasource configuration..."
cat > monitoring/grafana/provisioning/datasources/datasource.yml << EOF
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
EOF

echo "Creating Grafana dashboard configuration..."
cat > monitoring/grafana/provisioning/dashboards/dashboard.yml << EOF
apiVersion: 1
providers:
  - name: 'default'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    options:
      path: /var/lib/grafana/dashboards
EOF

echo "Building and starting all services..."
docker-compose -f docker-compose-all.yml up --build -d

echo "Waiting for services to start..."
sleep 30

echo "Checking service status..."
docker-compose -f docker-compose-all.yml ps

echo "==================================================="
echo "SparkFund Services:"
echo "==================================================="
echo "API Gateway:         http://localhost:8080"
echo "KYC Service:         http://localhost:8081"
echo "Investment Service:  http://localhost:8082"
echo "User Service:        http://localhost:8083"
echo "AI Service:          http://localhost:8001"
echo "==================================================="
echo "Monitoring:"
echo "==================================================="
echo "Prometheus:          http://localhost:9090"
echo "Grafana:             http://localhost:3000 (admin/admin)"
echo "Jaeger:              http://localhost:16686"
echo "==================================================="

echo "To view logs, run: docker-compose -f docker-compose-all.yml logs -f"
echo "To stop all services, run: docker-compose -f docker-compose-all.yml down"
