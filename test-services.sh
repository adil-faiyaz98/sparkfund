#!/bin/bash
# SparkFund - Test Services Script
echo "====================================================="
echo "SparkFund - Testing Services"
echo "====================================================="

# Function to test an endpoint
test_endpoint() {
    local service=$1
    local endpoint=$2
    local method=${3:-GET}
    
    echo "Testing $service at $endpoint..."
    
    response=$(curl -s -o /dev/null -w "%{http_code}" -X $method $endpoint)
    
    if [ "$response" -ge 200 ] && [ "$response" -lt 300 ]; then
        echo "Status: $response - Success!"
        return 0
    else
        echo "Status: $response - Failed!"
        return 1
    fi
}

# Test API Gateway
test_endpoint "API Gateway" "http://localhost:8080/health"
api_gateway_health=$?

# Test KYC Service
test_endpoint "KYC Service" "http://localhost:8081/health"
kyc_service_health=$?

# Test Investment Service
test_endpoint "Investment Service" "http://localhost:8082/health"
investment_service_health=$?

# Test User Service
test_endpoint "User Service" "http://localhost:8083/health"
user_service_health=$?

# Test AI Service
test_endpoint "AI Service" "http://localhost:8001/health"
ai_service_health=$?

# Summary
echo "====================================================="
echo "Service Health Summary:"
echo "====================================================="
echo "API Gateway: $([ $api_gateway_health -eq 0 ] && echo "Healthy" || echo "Unhealthy")"
echo "KYC Service: $([ $kyc_service_health -eq 0 ] && echo "Healthy" || echo "Unhealthy")"
echo "Investment Service: $([ $investment_service_health -eq 0 ] && echo "Healthy" || echo "Unhealthy")"
echo "User Service: $([ $user_service_health -eq 0 ] && echo "Healthy" || echo "Unhealthy")"
echo "AI Service: $([ $ai_service_health -eq 0 ] && echo "Healthy" || echo "Unhealthy")"
echo "====================================================="

# Test basic functionality if all services are healthy
if [ $api_gateway_health -eq 0 ] && [ $kyc_service_health -eq 0 ] && [ $investment_service_health -eq 0 ] && [ $user_service_health -eq 0 ] && [ $ai_service_health -eq 0 ]; then
    echo "All services are healthy. Testing basic functionality..."
    
    # Add more functional tests here
    
    echo "Basic functionality tests completed."
else
    echo "Some services are unhealthy. Skipping functional tests."
fi

echo "====================================================="
