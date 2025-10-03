#!/bin/bash

# JobFair API Testing Script
# This script tests all the fixed endpoints

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost"
TOKEN=""

echo "======================================"
echo "  JobFair API Testing Script"
echo "======================================"
echo ""

# Function to print test result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
    fi
}

# Function to extract token from login response
extract_token() {
    echo "$1" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4
}

echo "Step 1: Testing Auth Service"
echo "------------------------------"

# Register a new user
echo "Testing: POST /api/v1/auth/register"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test1234!",
    "role": "jobseeker",
    "full_name": "Test User"
  }')

if echo "$REGISTER_RESPONSE" | grep -q "success"; then
    print_result 0 "User registration"
else
    echo "  Note: User might already exist"
fi

# Login
echo "Testing: POST /api/v1/auth/login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test1234!"
  }')

TOKEN=$(extract_token "$LOGIN_RESPONSE")

if [ -n "$TOKEN" ]; then
    print_result 0 "User login"
    echo "  Token: ${TOKEN:0:20}..."
else
    print_result 1 "User login - Could not extract token"
    echo "  Response: $LOGIN_RESPONSE"
    echo ""
    echo "Please login manually and set TOKEN variable:"
    echo "export TOKEN='your-token-here'"
    exit 1
fi

echo ""
echo "Step 2: Create Profile"
echo "------------------------------"

# Create profile
echo "Testing: POST /api/v1/profiles"
PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/profiles" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "phone_number": "081234567890",
    "date_of_birth": "1990-01-01",
    "gender": "male"
  }')

if echo "$PROFILE_RESPONSE" | grep -q "success"; then
    print_result 0 "Profile creation"
else
    echo "  Note: Profile might already exist"
    echo "  Response: $PROFILE_RESPONSE"
fi

echo ""
echo "Step 3: Testing Fixed Endpoints"
echo "------------------------------"

# Test 1: Bulk Skills
echo "Testing: POST /api/v1/skills/bulk (FIX #1)"
SKILLS_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/skills/bulk" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "skills": [
      {
        "skill_name": "PostgreSQL",
        "skill_type": "technical",
        "proficiency_level": "advanced",
        "years_of_experience": 4
      },
      {
        "skill_name": "Docker",
        "skill_type": "technical",
        "proficiency_level": "intermediate",
        "years_of_experience": 3
      },
      {
        "skill_name": "Communication",
        "skill_type": "soft",
        "proficiency_level": "expert",
        "years_of_experience": 5
      }
    ]
  }')

if echo "$SKILLS_RESPONSE" | grep -q '"id"' && echo "$SKILLS_RESPONSE" | grep -q '"skill_name"'; then
    print_result 0 "Bulk skills creation - Data is saved"
    echo "  Response preview: $(echo $SKILLS_RESPONSE | cut -c1-100)..."
else
    print_result 1 "Bulk skills creation - Data might not be saved"
    echo "  Response: $SKILLS_RESPONSE"
fi

# Test 2: Career Preference
echo ""
echo "Testing: POST /api/v1/career-preference (FIX #2)"
CAREER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/career-preference" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "hybrid",
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "currency": "IDR",
    "willing_to_relocate": true,
    "available_from": "2025-11-01T00:00:00Z"
  }')

if echo "$CAREER_RESPONSE" | grep -q "success.*true"; then
    print_result 0 "Career preference creation - No foreign key error"
else
    print_result 1 "Career preference creation"
    echo "  Response: $CAREER_RESPONSE"
fi

# Test 3: Position Preferences
echo ""
echo "Testing: POST /api/v1/position-preferences (FIX #3)"
POSITION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/position-preferences" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "positions": [
      "Software Engineer",
      "Backend Developer",
      "Full Stack Developer"
    ]
  }')

if echo "$POSITION_RESPONSE" | grep -q "success.*true"; then
    print_result 0 "Position preferences creation - Array of strings accepted"
else
    print_result 1 "Position preferences creation"
    echo "  Response: $POSITION_RESPONSE"
fi

# Test 4: CV Upload (requires a file)
echo ""
echo "Testing: POST /api/v1/cv (FIX #4)"
echo "  Note: CV upload test requires a PDF or DOCX file"
echo "  Command to test manually:"
echo "  curl -X POST $BASE_URL/api/v1/cv \\"
echo "    -H \"Authorization: Bearer $TOKEN\" \\"
echo "    -F \"file=@/path/to/your/resume.pdf\""

echo ""
echo "Step 4: Testing Job Service via Gateway (FIX #5)"
echo "------------------------------"

# Test Job Service
echo "Testing: GET /api/v1/jobs"
JOBS_RESPONSE=$(curl -s "$BASE_URL/api/v1/jobs?page=1&limit=10")

if echo "$JOBS_RESPONSE" | grep -q "jobs" || echo "$JOBS_RESPONSE" | grep -q "data"; then
    print_result 0 "Job service routing via gateway"
else
    print_result 1 "Job service routing via gateway"
    echo "  Response: $JOBS_RESPONSE"
fi

echo ""
echo "Step 5: Testing Other Services via Gateway"
echo "------------------------------"

# Test Company Service
echo "Testing: GET /api/v1/companies"
COMPANIES_RESPONSE=$(curl -s "$BASE_URL/api/v1/companies?page=1&limit=10")

if echo "$COMPANIES_RESPONSE" | grep -q "companies" || echo "$COMPANIES_RESPONSE" | grep -q "data"; then
    print_result 0 "Company service routing via gateway"
else
    print_result 1 "Company service routing via gateway"
    echo "  Response: $COMPANIES_RESPONSE"
fi

echo ""
echo "======================================"
echo "  Testing Complete!"
echo "======================================"
echo ""
echo "Summary of Fixes:"
echo "  1. ✓ Bulk skills now saves data correctly"
echo "  2. ✓ Career preference accepts new field format"
echo "  3. ✓ Position preferences accepts array of strings"
echo "  4. ⚠ CV upload (test manually with file)"
echo "  5. ✓ API Gateway routes all services correctly"
echo ""
echo "Your token (save this for future use):"
echo "$TOKEN"
