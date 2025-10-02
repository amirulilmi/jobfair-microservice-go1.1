#!/bin/bash

# JobFair User Profile Service - API Testing Script
# This script tests all endpoints of the User Profile Service

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8083"
AUTH_SERVICE_URL="http://localhost:8080"

# Test counter
PASSED=0
FAILED=0
TOTAL=0

# Function to print test result
print_result() {
    local test_name=$1
    local status=$2
    local response=$3
    
    TOTAL=$((TOTAL + 1))
    
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}✓${NC} Test $TOTAL: $test_name ${GREEN}PASSED${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗${NC} Test $TOTAL: $test_name ${RED}FAILED${NC}"
        echo -e "${RED}Response: $response${NC}"
        FAILED=$((FAILED + 1))
    fi
}

# Function to make API call
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    if [ -n "$token" ]; then
        if [ -n "$data" ]; then
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Authorization: Bearer $token" \
                -H "Content-Type: application/json" \
                -d "$data"
        else
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Authorization: Bearer $token"
        fi
    else
        curl -s -X $method "$BASE_URL$endpoint"
    fi
}

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}JobFair User Profile Service - API Testing${NC}"
echo -e "${BLUE}============================================${NC}\n"

# ============================================
# TEST 1: Health Check (No Auth)
# ============================================
echo -e "\n${YELLOW}[1] Testing Health Check...${NC}"
response=$(curl -s "$BASE_URL/health")
if echo "$response" | grep -q "ok"; then
    print_result "Health Check" 0 "$response"
else
    print_result "Health Check" 1 "$response"
fi

# ============================================
# SETUP: Get JWT Token from Auth Service
# ============================================
echo -e "\n${YELLOW}[SETUP] Getting JWT Token from Auth Service...${NC}"
echo "Please provide authentication:"
read -p "Enter email: " EMAIL
read -sp "Enter password: " PASSWORD
echo ""

login_response=$(curl -s -X POST "$AUTH_SERVICE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo $login_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Failed to get JWT token. Please make sure:${NC}"
    echo -e "${RED}  1. Auth service is running at $AUTH_SERVICE_URL${NC}"
    echo -e "${RED}  2. Email and password are correct${NC}"
    echo -e "${RED}  3. User is registered${NC}"
    exit 1
fi

echo -e "${GREEN}✓ JWT Token obtained successfully${NC}"
echo -e "${BLUE}Token: ${TOKEN:0:20}...${NC}\n"

# ============================================
# TEST 2: Create Profile
# ============================================
echo -e "\n${YELLOW}[2] Testing Create Profile...${NC}"
response=$(api_call "POST" "/api/v1/profiles" \
    '{"full_name":"Test User","phone_number":"081234567890"}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Profile" 0 "$response"
else
    print_result "Create Profile" 1 "$response"
fi

# ============================================
# TEST 3: Get Profile
# ============================================
echo -e "\n${YELLOW}[3] Testing Get Profile...${NC}"
response=$(api_call "GET" "/api/v1/profiles" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get Profile" 0 "$response"
else
    print_result "Get Profile" 1 "$response"
fi

# ============================================
# TEST 4: Update Profile
# ============================================
echo -e "\n${YELLOW}[4] Testing Update Profile...${NC}"
response=$(api_call "PUT" "/api/v1/profiles" \
    '{"bio":"Software Engineer with passion for microservices","city":"Jakarta"}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Update Profile" 0 "$response"
else
    print_result "Update Profile" 1 "$response"
fi

# ============================================
# TEST 5: Create Work Experience
# ============================================
echo -e "\n${YELLOW}[5] Testing Create Work Experience...${NC}"
response=$(api_call "POST" "/api/v1/work-experiences" \
    '{"company_name":"PT Test","job_position":"Software Engineer","start_date":"2020-01-01T00:00:00Z","is_current_job":true,"job_description":"Developing microservices"}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Work Experience" 0 "$response"
    WORK_EXP_ID=$(echo $response | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
else
    print_result "Create Work Experience" 1 "$response"
fi

# ============================================
# TEST 6: Get All Work Experiences
# ============================================
echo -e "\n${YELLOW}[6] Testing Get All Work Experiences...${NC}"
response=$(api_call "GET" "/api/v1/work-experiences" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get All Work Experiences" 0 "$response"
else
    print_result "Get All Work Experiences" 1 "$response"
fi

# ============================================
# TEST 7: Create Education
# ============================================
echo -e "\n${YELLOW}[7] Testing Create Education...${NC}"
response=$(api_call "POST" "/api/v1/educations" \
    '{"university":"Universitas Indonesia","major":"Computer Science","degree":"Bachelor","start_date":"2016-08-01T00:00:00Z","end_date":"2020-07-31T00:00:00Z","gpa":3.75}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Education" 0 "$response"
else
    print_result "Create Education" 1 "$response"
fi

# ============================================
# TEST 8: Get All Educations
# ============================================
echo -e "\n${YELLOW}[8] Testing Get All Educations...${NC}"
response=$(api_call "GET" "/api/v1/educations" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get All Educations" 0 "$response"
else
    print_result "Get All Educations" 1 "$response"
fi

# ============================================
# TEST 9: Create Certification
# ============================================
echo -e "\n${YELLOW}[9] Testing Create Certification...${NC}"
response=$(api_call "POST" "/api/v1/certifications" \
    '{"certification_name":"AWS Certified Developer","organizer":"Amazon Web Services","issue_date":"2023-01-15T00:00:00Z"}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Certification" 0 "$response"
else
    print_result "Create Certification" 1 "$response"
fi

# ============================================
# TEST 10: Get All Certifications
# ============================================
echo -e "\n${YELLOW}[10] Testing Get All Certifications...${NC}"
response=$(api_call "GET" "/api/v1/certifications" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get All Certifications" 0 "$response"
else
    print_result "Get All Certifications" 1 "$response"
fi

# ============================================
# TEST 11: Create Skills (Bulk)
# ============================================
echo -e "\n${YELLOW}[11] Testing Create Skills (Bulk)...${NC}"
response=$(api_call "POST" "/api/v1/skills/bulk" \
    '{"technical_skills":[{"skill_name":"Go","skill_type":"technical","proficiency_level":"expert"},{"skill_name":"PostgreSQL","skill_type":"technical","proficiency_level":"advanced"}],"soft_skills":[{"skill_name":"Leadership","skill_type":"soft","proficiency_level":"advanced"}]}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Skills (Bulk)" 0 "$response"
else
    print_result "Create Skills (Bulk)" 1 "$response"
fi

# ============================================
# TEST 12: Get All Skills
# ============================================
echo -e "\n${YELLOW}[12] Testing Get All Skills...${NC}"
response=$(api_call "GET" "/api/v1/skills" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get All Skills" 0 "$response"
else
    print_result "Get All Skills" 1 "$response"
fi

# ============================================
# TEST 13: Create Career Preference
# ============================================
echo -e "\n${YELLOW}[13] Testing Create Career Preference...${NC}"
response=$(api_call "POST" "/api/v1/career-preference" \
    '{"is_actively_looking":true,"expected_salary_min":15000000,"expected_salary_max":25000000,"salary_currency":"IDR","is_negotiable":true,"preferred_work_types":["remote","hybrid"],"preferred_locations":["Jakarta"],"willing_to_relocate":false}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Career Preference" 0 "$response"
else
    print_result "Create Career Preference" 1 "$response"
fi

# ============================================
# TEST 14: Get Career Preference
# ============================================
echo -e "\n${YELLOW}[14] Testing Get Career Preference...${NC}"
response=$(api_call "GET" "/api/v1/career-preference" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get Career Preference" 0 "$response"
else
    print_result "Get Career Preference" 1 "$response"
fi

# ============================================
# TEST 15: Create Position Preferences
# ============================================
echo -e "\n${YELLOW}[15] Testing Create Position Preferences...${NC}"
response=$(api_call "POST" "/api/v1/position-preferences" \
    '{"positions":[{"position_name":"Software Engineer","priority":1},{"position_name":"Backend Developer","priority":2}]}' \
    "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Create Position Preferences" 0 "$response"
else
    print_result "Create Position Preferences" 1 "$response"
fi

# ============================================
# TEST 16: Get Position Preferences
# ============================================
echo -e "\n${YELLOW}[16] Testing Get Position Preferences...${NC}"
response=$(api_call "GET" "/api/v1/position-preferences" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get Position Preferences" 0 "$response"
else
    print_result "Get Position Preferences" 1 "$response"
fi

# ============================================
# TEST 17: Get Complete Profile
# ============================================
echo -e "\n${YELLOW}[17] Testing Get Complete Profile...${NC}"
response=$(api_call "GET" "/api/v1/profiles/full" "" "$TOKEN")

if echo "$response" | grep -q "success"; then
    print_result "Get Complete Profile" 0 "$response"
else
    print_result "Get Complete Profile" 1 "$response"
fi

# ============================================
# TEST 18: Get Profile Completion Status
# ============================================
echo -e "\n${YELLOW}[18] Testing Get Profile Completion Status...${NC}"
response=$(api_call "GET" "/api/v1/profiles/completion" "" "$TOKEN")

if echo "$response" | grep -q "completion_status"; then
    print_result "Get Profile Completion Status" 0 "$response"
    COMPLETION=$(echo $response | grep -o '"completion_status":[0-9]*' | cut -d':' -f2)
    echo -e "${BLUE}Profile Completion: $COMPLETION%${NC}"
else
    print_result "Get Profile Completion Status" 1 "$response"
fi

# ============================================
# SUMMARY
# ============================================
echo -e "\n${BLUE}============================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}============================================${NC}"
echo -e "Total Tests: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed! Service is working correctly.${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️  Some tests failed. Please check the errors above.${NC}"
    exit 1
fi
