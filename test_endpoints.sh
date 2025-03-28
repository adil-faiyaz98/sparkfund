#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Test JWT Token (replace with your generated token)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXItMTIzIiwiZXhwIjoxNzExNjY0ODAwLCJpYXQiOjE3MTE1Nzg0MDB9.YourGeneratedTokenHere"

echo "Testing API Gateway endpoints..."

# Health check
echo -e "\n${GREEN}Testing health endpoint...${NC}"
curl -s http://localhost:8080/health

# Metrics
echo -e "\n${GREEN}Testing metrics endpoint...${NC}"
curl -s http://localhost:8080/metrics

# Investments
echo -e "\n${GREEN}Testing investments endpoints...${NC}"

# List investments
echo -e "\nListing investments..."
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/investments

# Create investment
echo -e "\nCreating investment..."
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test-user-123",
    "amount": 1000,
    "type": "STOCK",
    "symbol": "MSFT",
    "quantity": 10,
    "purchase_price": 280.50
  }' \
  http://localhost:8080/api/v1/investments

# Portfolios
echo -e "\n${GREEN}Testing portfolio endpoints...${NC}"

# List portfolios
echo -e "\nListing portfolios..."
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/portfolios

# Create portfolio
echo -e "\nCreating portfolio..."
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test-user-123",
    "name": "Test Portfolio",
    "description": "A test portfolio"
  }' \
  http://localhost:8080/api/v1/portfolios

echo -e "\n${GREEN}Testing completed!${NC}" 