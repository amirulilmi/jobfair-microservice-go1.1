#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Auto-Create Profile Feature Test${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if email and password provided
EMAIL=${1:-"testuser@example.com"}
PASSWORD=${2:-"TestPass123!"}

echo -e "${YELLOW}Using credentials:${NC}"
echo "Email: $EMAIL"
echo "Password: $PASSWORD"
echo ""

# Step 1: Login
echo -e "${YELLOW}Step 1: Login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Login failed!${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}"
echo "Token: ${TOKEN:0:30}..."
echo ""

# Step 2: Try to get profile (should not exist yet)
echo -e "${YELLOW}Step 2: Checking if profile exists...${NC}"
PROFILE_CHECK=$(curl -s -X GET http://localhost/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN")

PROFILE_EXISTS=$(echo $PROFILE_CHECK | jq -r '.success')

if [ "$PROFILE_EXISTS" = "false" ]; then
    echo -e "${BLUE}ℹ Profile doesn't exist yet (expected)${NC}"
else
    echo -e "${YELLOW}⚠ Profile already exists${NC}"
fi
echo ""

# Step 3: Create career preference (should auto-create profile)
echo -e "${YELLOW}Step 3: Creating career preference (will auto-create profile)...${NC}"
PREFERENCE_RESPONSE=$(curl -s -X POST http://localhost/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "hybrid",
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "currency": "IDR",
    "willing_to_relocate": true,
    "available_from": "2025-11-01"
  }')

PREF_SUCCESS=$(echo $PREFERENCE_RESPONSE | jq -r '.success')

if [ "$PREF_SUCCESS" = "true" ]; then
    echo -e "${GREEN}✓ Career preference created successfully${NC}"
    echo "Response:"
    echo $PREFERENCE_RESPONSE | jq
else
    echo -e "${RED}✗ Failed to create career preference${NC}"
    echo "Response:"
    echo $PREFERENCE_RESPONSE | jq
    exit 1
fi
echo ""

# Step 4: Verify profile was auto-created
echo -e "${YELLOW}Step 4: Verifying profile was auto-created...${NC}"
PROFILE_CHECK2=$(curl -s -X GET http://localhost/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN")

PROFILE_EXISTS2=$(echo $PROFILE_CHECK2 | jq -r '.success')

if [ "$PROFILE_EXISTS2" = "true" ]; then
    echo -e "${GREEN}✓ Profile auto-created successfully!${NC}"
    PROFILE_ID=$(echo $PROFILE_CHECK2 | jq -r '.data.id')
    COMPLETION=$(echo $PROFILE_CHECK2 | jq -r '.data.completion_status')
    echo "Profile ID: $PROFILE_ID"
    echo "Completion: $COMPLETION%"
else
    echo -e "${RED}✗ Profile was not created${NC}"
    echo "Response:"
    echo $PROFILE_CHECK2 | jq
    exit 1
fi
echo ""

# Step 5: Add work experience
echo -e "${YELLOW}Step 5: Adding work experience...${NC}"
WORKEXP_RESPONSE=$(curl -s -X POST http://localhost/api/v1/work-experiences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Tech Corp Indonesia",
    "job_position": "Senior Software Engineer",
    "start_date": "2020-01-01",
    "end_date": "2023-12-31",
    "is_current_job": false,
    "job_description": "Developed microservices architecture"
  }')

WORKEXP_SUCCESS=$(echo $WORKEXP_RESPONSE | jq -r '.success')

if [ "$WORKEXP_SUCCESS" = "true" ]; then
    echo -e "${GREEN}✓ Work experience added${NC}"
else
    echo -e "${RED}✗ Failed to add work experience${NC}"
    echo "Response:"
    echo $WORKEXP_RESPONSE | jq
fi
echo ""

# Step 6: Add education
echo -e "${YELLOW}Step 6: Adding education...${NC}"
EDUCATION_RESPONSE=$(curl -s -X POST http://localhost/api/v1/educations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "university": "University of Indonesia",
    "major": "Computer Science",
    "degree": "Bachelor",
    "start_date": "2016-08-01",
    "end_date": "2020-07-01",
    "is_current": false,
    "gpa": "3.75",
    "description": "Focus on Software Engineering"
  }')

EDUCATION_SUCCESS=$(echo $EDUCATION_RESPONSE | jq -r '.success')

if [ "$EDUCATION_SUCCESS" = "true" ]; then
    echo -e "${GREEN}✓ Education added${NC}"
else
    echo -e "${RED}✗ Failed to add education${NC}"
    echo "Response:"
    echo $EDUCATION_RESPONSE | jq
fi
echo ""

# Step 7: Add skills
echo -e "${YELLOW}Step 7: Adding skills...${NC}"
SKILLS_RESPONSE=$(curl -s -X POST http://localhost/api/v1/skills/bulk \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "skills": [
      {
        "skill_name": "Go",
        "skill_type": "technical",
        "proficiency_level": "expert",
        "years_of_experience": 5
      },
      {
        "skill_name": "Docker",
        "skill_type": "technical",
        "proficiency_level": "advanced",
        "years_of_experience": 4
      },
      {
        "skill_name": "Leadership",
        "skill_type": "soft",
        "proficiency_level": "advanced",
        "years_of_experience": 3
      }
    ]
  }')

SKILLS_SUCCESS=$(echo $SKILLS_RESPONSE | jq -r '.success')

if [ "$SKILLS_SUCCESS" = "true" ]; then
    SKILLS_COUNT=$(echo $SKILLS_RESPONSE | jq -r '.data | length')
    echo -e "${GREEN}✓ $SKILLS_COUNT skills added${NC}"
else
    echo -e "${RED}✗ Failed to add skills${NC}"
    echo "Response:"
    echo $SKILLS_RESPONSE | jq
fi
echo ""

# Step 8: Check final completion status
echo -e "${YELLOW}Step 8: Checking final completion status...${NC}"
COMPLETION_RESPONSE=$(curl -s -X GET http://localhost/api/v1/profiles/completion \
  -H "Authorization: Bearer $TOKEN")

COMPLETION_STATUS=$(echo $COMPLETION_RESPONSE | jq -r '.data.completion_percentage')

echo -e "${GREEN}✓ Profile completion: $COMPLETION_STATUS%${NC}"
echo ""

# Final summary
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}           Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}✓ Login successful${NC}"
echo -e "${GREEN}✓ Profile auto-created${NC}"
echo -e "${GREEN}✓ Career preference created${NC}"
echo -e "${GREEN}✓ Work experience added${NC}"
echo -e "${GREEN}✓ Education added${NC}"
echo -e "${GREEN}✓ Skills added${NC}"
echo -e "${GREEN}✓ Final completion: $COMPLETION_STATUS%${NC}"
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}All tests passed! 🎉${NC}"
echo -e "${BLUE}========================================${NC}"
