#!/bin/bash

# Test script untuk verify semua bug fixes
# Usage: ./scripts/test-bug-fixes.sh [TOKEN]

set -e

# Configuration
API_GATEWAY="http://localhost:8000"
TOKEN="${1:-}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🧪 Testing Bug Fixes"
echo "==================="
echo ""

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⚠️  No token provided. Some tests will be skipped.${NC}"
    echo "Usage: $0 YOUR_JWT_TOKEN"
    echo ""
fi

# Test 1: Profile Creation (with headline, summary, etc.)
echo -e "${BLUE}Test 1: Profile Creation with New Fields${NC}"
if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⊘ Skipped (no token)${NC}"
else
    RESPONSE=$(curl -s -X POST "$API_GATEWAY/api/v1/profiles" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "first_name": "Test",
            "last_name": "User",
            "headline": "Software Engineer",
            "summary": "Test summary",
            "location": "Jakarta, Indonesia",
            "phone_number": "081234567890",
            "github_url": "https://github.com/testuser"
        }')
    
    if echo "$RESPONSE" | grep -q '"success":true'; then
        echo -e "${GREEN}✓ Profile created successfully${NC}"
        
        # Check if full_name was auto-generated
        if echo "$RESPONSE" | grep -q '"full_name":"Test User"'; then
            echo -e "${GREEN}✓ Full name auto-generated correctly${NC}"
        fi
        
        # Check new fields
        if echo "$RESPONSE" | grep -q '"headline":"Software Engineer"'; then
            echo -e "${GREEN}✓ Headline field working${NC}"
        fi
        
        if echo "$RESPONSE" | grep -q '"github_url"'; then
            echo -e "${GREEN}✓ GitHub URL field working${NC}"
        fi
    else
        echo -e "${RED}✗ Profile creation failed${NC}"
        echo "Response: $RESPONSE"
    fi
fi
echo ""

# Test 2: Profile Update (first_name only)
echo -e "${BLUE}Test 2: Profile Update with First Name Only${NC}"
if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⊘ Skipped (no token)${NC}"
else
    RESPONSE=$(curl -s -X PUT "$API_GATEWAY/api/v1/profiles" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "first_name": "UpdatedName"
        }')
    
    if echo "$RESPONSE" | grep -q '"success":true'; then
        echo -e "${GREEN}✓ Profile updated with first_name${NC}"
        
        if echo "$RESPONSE" | grep -q '"first_name":"UpdatedName"'; then
            echo -e "${GREEN}✓ First name updated correctly${NC}"
        fi
    else
        echo -e "${RED}✗ Profile update failed${NC}"
        echo "Response: $RESPONSE"
    fi
fi
echo ""

# Test 3: Valid JSON (Job Application)
echo -e "${BLUE}Test 3: Job Application with Valid JSON${NC}"
if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⊘ Skipped (no token)${NC}"
else
    # First, get a job ID
    JOBS_RESPONSE=$(curl -s -X GET "$API_GATEWAY/api/v1/jobs?limit=1")
    JOB_ID=$(echo "$JOBS_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    
    if [ -n "$JOB_ID" ]; then
        echo "Testing with job ID: $JOB_ID"
        
        # Valid JSON (no comments)
        RESPONSE=$(curl -s -X POST "$API_GATEWAY/api/v1/jobs/$JOB_ID/apply" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d '{
                "cv_url": "https://example.com/cv.pdf",
                "cover_letter": "I am interested in this position"
            }')
        
        if echo "$RESPONSE" | grep -q '"success":true'; then
            echo -e "${GREEN}✓ Job application successful${NC}"
        elif echo "$RESPONSE" | grep -q "already applied"; then
            echo -e "${YELLOW}⚠ Already applied to this job (expected for repeated tests)${NC}"
        elif echo "$RESPONSE" | grep -q "invalid character"; then
            echo -e "${RED}✗ JSON parsing still failing${NC}"
            echo "Response: $RESPONSE"
        else
            echo -e "${YELLOW}⚠ Other error (may be expected)${NC}"
            echo "Response: $RESPONSE"
        fi
    else
        echo -e "${YELLOW}⊘ No jobs available for testing${NC}"
    fi
fi
echo ""

# Test 4: Database Schema Check
echo -e "${BLUE}Test 4: Database Schema Verification${NC}"
COLUMNS_CHECK=$(docker exec -i postgres-profile psql -U jobfair_user -d jobfair_profiles -t -c "
    SELECT column_name 
    FROM information_schema.columns 
    WHERE table_name = 'profiles' 
    AND column_name IN ('first_name', 'last_name', 'headline', 'summary', 'location', 'github_url')
    ORDER BY column_name;
" 2>&1)

EXPECTED_COLS=("first_name" "github_url" "headline" "last_name" "location" "summary")
for col in "${EXPECTED_COLS[@]}"; do
    if echo "$COLUMNS_CHECK" | grep -q "$col"; then
        echo -e "  ${GREEN}✓${NC} Column '$col' exists"
    else
        echo -e "  ${RED}✗${NC} Column '$col' missing"
    fi
done
echo ""

# Test 5: API Documentation Check
echo -e "${BLUE}Test 5: Documentation Files${NC}"
if [ -f "BUG_FIXES_OCT_2025.md" ]; then
    echo -e "${GREEN}✓ Bug fixes documentation exists${NC}"
else
    echo -e "${RED}✗ Bug fixes documentation missing${NC}"
fi

if [ -f "scripts/apply-bug-fixes.sh" ]; then
    echo -e "${GREEN}✓ Fix application script exists${NC}"
else
    echo -e "${RED}✗ Fix application script missing${NC}"
fi
echo ""

# Summary
echo "==================="
echo -e "${GREEN}Testing Complete${NC}"
echo "==================="
echo ""
echo "Tips for manual testing:"
echo ""
echo "1. Company Logo Upload:"
echo "   curl -X POST $API_GATEWAY/api/v1/companies/1/logo \\"
echo "     -H \"Authorization: Bearer \$TOKEN\" \\"
echo "     -F \"file=@logo.png\""
echo ""
echo "2. Job Application (ensure NO comments in JSON):"
echo "   curl -X POST $API_GATEWAY/api/v1/jobs/1/apply \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -H \"Authorization: Bearer \$TOKEN\" \\"
echo "     -d '{\"cv_url\":\"...\",\"cover_letter\":\"...\"}'"
echo ""
echo "3. Profile with first_name + last_name:"
echo "   curl -X POST $API_GATEWAY/api/v1/profiles \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -H \"Authorization: Bearer \$TOKEN\" \\"
echo "     -d '{\"first_name\":\"John\",\"last_name\":\"Doe\"}'"
echo ""
