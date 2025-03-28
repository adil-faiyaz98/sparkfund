#!/bin/bash

# Create new directory structure
mkdir -p sparkfund/{.github/workflows,api/{openapi,proto},build/{ci,package},config,deployments/{docker,kubernetes/helm},docs/{architecture,api,development},pkg/{logger,metrics,validator},scripts,services}

# Move existing services to services directory
mv sparkfund/aml-service sparkfund/services/
mv sparkfund/auth-service sparkfund/services/
mv sparkfund/investment-service sparkfund/services/
mv sparkfund/kyc-service sparkfund/services/
mv sparkfund/credit-scoring-service sparkfund/services/

# Move Helm charts to deployments
mv sparkfund/helm/* sparkfund/deployments/kubernetes/helm/

# Move observability components to deployments
mv sparkfund/{jaeger,grafana,prometheus} sparkfund/deployments/kubernetes/

# Move configuration
mv sparkfund/config/* sparkfund/config/

# Move scripts
mv sparkfund/scripts/* sparkfund/scripts/

# Move GitHub Actions workflows
mv sparkfund/.github/workflows/* sparkfund/.github/workflows/

# Move development docker-compose
mv sparkfund/docker-compose.dev.yml sparkfund/deployments/docker/

# Create service-specific directories for each service
for service in aml-service auth-service investment-service kyc-service credit-scoring-service; do
    mkdir -p sparkfund/services/$service/{cmd/server,internal/{domain,infrastructure,interfaces,usecase},configs,deployments,docs,test}
done

# Create root go.mod if it doesn't exist
if [ ! -f sparkfund/go.mod ]; then
    cat > sparkfund/go.mod << EOF
module github.com/sparkfund

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v8 v8.11.5
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/google/uuid v1.3.0
    github.com/joho/godotenv v1.5.1
    github.com/prometheus/client_golang v1.16.0
    github.com/sirupsen/logrus v1.9.3
    gorm.io/driver/postgres v1.4.7
    gorm.io/gorm v1.25.5
)
EOF
fi

echo "Project structure reorganized successfully!" 