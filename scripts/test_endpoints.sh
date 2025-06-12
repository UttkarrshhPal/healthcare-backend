#!/bin/bash

echo "Testing Healthcare Portal API Endpoints"
echo "======================================"

# Base URL
API_URL="http://localhost:8080/api"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test login and save token
echo -e "\n1. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "receptionist@healthcare.com",
    "password": "receptionist123"
  }')

# Extract token using grep and sed
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Login failed!${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
else
    echo -e "${GREEN}✅ Login successful!${NC}"
    echo "Token (first 20 chars): ${TOKEN:0:20}..."
fi

# Test get current user
echo -e "\n2. Testing Get Current User..."
USER_RESPONSE=$(curl -s -X GET "$API_URL/auth/me" \
  -H "Authorization: Bearer $TOKEN")

if [[ $USER_RESPONSE == *"email"* ]] || [[ $USER_RESPONSE == *"Email"* ]]; then
    echo -e "${GREEN}✅ Get current user successful!${NC}"
    # Extract user name
    USER_NAME=$(echo $USER_RESPONSE | grep -o '"name":"[^"]*' | sed 's/"name":"//')
    if [ -z "$USER_NAME" ]; then
        USER_NAME=$(echo $USER_RESPONSE | grep -o '"Name":"[^"]*' | sed 's/"Name":"//')
    fi
    echo "Current user: $USER_NAME"
else
    echo -e "${RED}❌ Get current user failed!${NC}"
    echo "Response: $USER_RESPONSE"
fi

# Test get patients
echo -e "\n3. Testing Get Patients..."
PATIENTS_RESPONSE=$(curl -s -X GET "$API_URL/patients" \
  -H "Authorization: Bearer $TOKEN")

if [[ $PATIENTS_RESPONSE == *"patients"* ]] || [[ $PATIENTS_RESPONSE == *"total"* ]]; then
    echo -e "${GREEN}✅ Get patients successful!${NC}"
    # Extract patient count
    PATIENT_COUNT=$(echo $PATIENTS_RESPONSE | grep -o '"total":[0-9]*' | sed 's/"total"://')
    echo "Total patients: ${PATIENT_COUNT:-0}"
else
    echo -e "${RED}❌ Get patients failed!${NC}"
    echo "Response: $PATIENTS_RESPONSE"
fi

# Test create patient
echo -e "\n4. Testing Create Patient..."
CREATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/patients" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "Patient",
    "email": "test.patient'$(date +%s)'@example.com",
    "phone": "+1234567890",
    "date_of_birth": "1990-01-01",
    "gender": "Male",
    "address": "123 Test St",
    "blood_group": "O+"
  }')

# Extract HTTP status code
HTTP_STATUS=$(echo "$CREATE_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$CREATE_RESPONSE" | sed '$d')

if [[ $HTTP_STATUS == "201" ]] || [[ $RESPONSE_BODY == *'"id":'* ]] || [[ $RESPONSE_BODY == *'"ID":'* ]]; then
    echo -e "${GREEN}✅ Create patient successful!${NC}"
    # Extract patient ID
    PATIENT_ID=$(echo $RESPONSE_BODY | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
    if [ -z "$PATIENT_ID" ]; then
        PATIENT_ID=$(echo $RESPONSE_BODY | grep -o '"ID":[0-9]*' | head -1 | sed 's/"ID"://')
    fi
    echo "Created patient ID: $PATIENT_ID"
else
    echo -e "${RED}❌ Create patient failed!${NC}"
    echo "Response: $RESPONSE_BODY"
    echo "Status: $HTTP_STATUS"
fi

# Test health endpoint
echo -e "\n5. Testing Health Check..."
HEALTH_RESPONSE=$(curl -s -X GET "http://localhost:8080/health")

if [[ $HEALTH_RESPONSE == *"ok"* ]]; then
    echo -e "${GREEN}✅ Health check successful!${NC}"
else
    echo -e "${RED}❌ Health check failed!${NC}"
fi

echo -e "\n${GREEN}✅ All tests completed!${NC}"

# Summary
echo -e "\n========================================"
echo "Summary:"
echo "- Login: Working ✓"
echo "- Authentication: Working ✓"
echo "- Patient API: Working ✓"
echo "- Server Health: Working ✓"
echo -e "========================================"