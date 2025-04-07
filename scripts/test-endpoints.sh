#!/bin/bash

# Shell script to test all available endpoints in SparkFund

echo "====================================================="
echo "SparkFund - Testing Available Endpoints"
echo "====================================================="
echo ""

# Define the JWT token for authentication
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxOTE2MjM5MDIyLCJyb2xlcyI6WyJhZG1pbiIsInVzZXIiXX0.Ks0I-dCdjWUxJEwuGP0qlyYJGXXjUYlLCRwPIZXI5Ss"
HEADERS="-H \"Content-Type: application/json\" -H \"Authorization: Bearer $TOKEN\""

# Test GET endpoint
test_get() {
    local name=$1
    local url=$2
    
    echo -n "Testing $name: $url "
    
    response=$(curl -s -o response.txt -w "%{http_code}" -X GET $url $HEADERS --max-time 5)
    
    if [ "$response" == "200" ]; then
        echo -e "\e[32m[OK]\e[0m"
        echo -e "  Response: \e[90m$(cat response.txt)\e[0m"
    else
        echo -e "\e[31m[FAIL] Status code: $response\e[0m"
    fi
    
    echo ""
}

# Test POST endpoint
test_post() {
    local name=$1
    local url=$2
    local body=$3
    
    echo -n "Testing $name: $url "
    
    response=$(curl -s -o response.txt -w "%{http_code}" -X POST $url $HEADERS -d "$body" --max-time 5)
    
    if [ "$response" == "200" ]; then
        echo -e "\e[32m[OK]\e[0m"
        echo -e "  Response: \e[90m$(cat response.txt)\e[0m"
    else
        echo -e "\e[31m[FAIL] Status code: $response\e[0m"
    fi
    
    echo ""
}

# API Gateway
test_get "API Gateway" "http://localhost:8080/health"

# KYC Service
test_get "KYC Health" "http://localhost:8080/api/kyc/health"
test_get "KYC Status" "http://localhost:8080/api/kyc/api/v1/kyc/status"
test_post "KYC Verify" "http://localhost:8080/api/kyc/api/v1/kyc/verify" "{}"

# Investment Service
test_get "Investment Health" "http://localhost:8080/api/investment/health"
test_get "Investment List" "http://localhost:8080/api/investment/api/v1/investments"
test_post "Investment Create" "http://localhost:8080/api/investment/api/v1/investments/create" '{"user_id": "123", "amount": 1000, "type": "STOCK", "symbol": "AAPL", "quantity": 10}'

# User Service
test_get "User Health" "http://localhost:8080/api/user/health"
test_get "User List" "http://localhost:8080/api/user/api/v1/users"
test_post "User Register" "http://localhost:8080/api/user/api/v1/users/register" '{"email": "test@example.com", "first_name": "Test", "last_name": "User", "password": "password123"}'

# AI Service
test_get "AI Health" "http://localhost:8080/api/ai/health"

echo "====================================================="
echo "Testing Complete"
echo "====================================================="

# Clean up
rm -f response.txt
