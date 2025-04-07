#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080/api/v1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to make API requests
function make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    echo -e "${GREEN}Making $method request to $endpoint${NC}"
    
    if [ -n "$data" ]; then
        echo "Request data: $data"
    fi
    
    local headers="-H 'Content-Type: application/json'"
    if [ -n "$token" ]; then
        headers="$headers -H 'Authorization: Bearer $token'"
    fi
    
    local cmd="curl -s -X $method $headers"
    if [ -n "$data" ]; then
        cmd="$cmd -d '$data'"
    fi
    cmd="$cmd $BASE_URL$endpoint"
    
    local response=$(eval $cmd)
    echo "Response: $response"
    echo ""
    
    # Return the response for further processing
    echo $response
}

# Test health endpoint
echo -e "${GREEN}Testing health endpoint${NC}"
curl -s http://localhost:8080/health
echo -e "\n"

# Login
echo -e "${GREEN}Testing login${NC}"
login_response=$(make_request "POST" "/auth/login" '{"email":"admin@example.com","password":"password123","device_info":{"ip_address":"192.168.1.1","user_agent":"Mozilla/5.0","device_type":"Desktop","os":"Windows","browser":"Chrome","location":"New York, USA"}}')

# Extract token from login response
token=$(echo $login_response | grep -o '"token":"[^"]*' | sed 's/"token":"//')
echo -e "${GREEN}Token: $token${NC}"

# Test AI models endpoint
echo -e "${GREEN}Testing AI models endpoint${NC}"
make_request "GET" "/ai/models" "" "$token"

# Upload a document
echo -e "${GREEN}Testing document upload${NC}"
document_response=$(curl -s -X POST -H "Authorization: Bearer $token" -F "file=@./test_data/passport.jpg" -F "type=PASSPORT" -F "name=passport.jpg" $BASE_URL/documents)
echo "Response: $document_response"
echo ""

# Extract document ID
document_id=$(echo $document_response | grep -o '"id":"[^"]*' | sed 's/"id":"//')
echo -e "${GREEN}Document ID: $document_id${NC}"

# Upload a selfie
echo -e "${GREEN}Testing selfie upload${NC}"
selfie_response=$(curl -s -X POST -H "Authorization: Bearer $token" -F "file=@./test_data/selfie.jpg" -F "type=SELFIE" -F "name=selfie.jpg" $BASE_URL/documents)
echo "Response: $selfie_response"
echo ""

# Extract selfie ID
selfie_id=$(echo $selfie_response | grep -o '"id":"[^"]*' | sed 's/"id":"//')
echo -e "${GREEN}Selfie ID: $selfie_id${NC}"

# Create a verification
echo -e "${GREEN}Testing verification creation${NC}"
verification_response=$(make_request "POST" "/verifications" '{"user_id":"'$(echo $login_response | grep -o '"id":"[^"]*' | sed 's/"id":"//g')'","kyc_id":"'$(uuidgen)'","document_id":"'$document_id'","method":"AI","status":"PENDING"}' "$token")

# Extract verification ID
verification_id=$(echo $verification_response | grep -o '"id":"[^"]*' | sed 's/"id":"//')
echo -e "${GREEN}Verification ID: $verification_id${NC}"

# Test document analysis
echo -e "${GREEN}Testing document analysis${NC}"
make_request "POST" "/ai/analyze-document" '{"document_id":"'$document_id'","verification_id":"'$verification_id'"}' "$token"

# Test face matching
echo -e "${GREEN}Testing face matching${NC}"
make_request "POST" "/ai/match-faces" '{"document_id":"'$document_id'","selfie_id":"'$selfie_id'","verification_id":"'$verification_id'"}' "$token"

# Test risk analysis
echo -e "${GREEN}Testing risk analysis${NC}"
make_request "POST" "/ai/analyze-risk" '{"user_id":"'$(echo $login_response | grep -o '"id":"[^"]*' | sed 's/"id":"//g')'","verification_id":"'$verification_id'","device_info":{"ip_address":"192.168.1.1","user_agent":"Mozilla/5.0","device_type":"Desktop","os":"Windows","browser":"Chrome","location":"New York, USA"}}' "$token"

# Test anomaly detection
echo -e "${GREEN}Testing anomaly detection${NC}"
make_request "POST" "/ai/detect-anomalies" '{"user_id":"'$(echo $login_response | grep -o '"id":"[^"]*' | sed 's/"id":"//g')'","verification_id":"'$verification_id'","device_info":{"ip_address":"192.168.1.1","user_agent":"Mozilla/5.0","device_type":"Desktop","os":"Windows","browser":"Chrome","location":"New York, USA"}}' "$token"

# Test document processing
echo -e "${GREEN}Testing document processing${NC}"
make_request "POST" "/ai/process-document" '{"document_id":"'$document_id'","selfie_id":"'$selfie_id'","verification_id":"'$verification_id'","device_info":{"ip_address":"192.168.1.1","user_agent":"Mozilla/5.0","device_type":"Desktop","os":"Windows","browser":"Chrome","location":"New York, USA"}}' "$token"

# Get verification status
echo -e "${GREEN}Testing verification status${NC}"
make_request "GET" "/verifications/$verification_id" "" "$token"

echo -e "${GREEN}All tests completed${NC}"
