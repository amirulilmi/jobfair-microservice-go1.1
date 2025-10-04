#!/bin/bash

# Safe Company Mapping Sync Script
# This script syncs company mappings from company-service to job-service
# Uses API endpoints instead of direct database access

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
API_GATEWAY="${API_GATEWAY:-http://localhost:8000}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"

echo -e "${BLUE}🔄 Company Mapping Sync Tool${NC}"
echo "================================"
echo ""

# Check if admin token is provided
if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}❌ Error: ADMIN_TOKEN environment variable is required${NC}"
    echo ""
    echo "Usage:"
    echo "  export ADMIN_TOKEN='your-admin-jwt-token'"
    echo "  ./scripts/sync-company-mappings-safe.sh"
    echo ""
    exit 1
fi

echo -e "${BLUE}Step 1: Checking data consistency...${NC}"

# Check current state
HEALTH_CHECK=$(curl -s -X GET "$API_GATEWAY/api/v1/admin/health/data-consistency" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json")

echo "$HEALTH_CHECK" | jq '.'

STATUS=$(echo "$HEALTH_CHECK" | jq -r '.status')
JOBS_WITHOUT_MAPPING=$(echo "$HEALTH_CHECK" | jq -r '.checks.jobs_without_mapping')

if [ "$STATUS" == "healthy" ]; then
    echo -e "${GREEN}✓ System is healthy, no sync needed${NC}"
    exit 0
fi

if [ "$JOBS_WITHOUT_MAPPING" -gt 0 ]; then
    echo -e "${YELLOW}⚠ Found $JOBS_WITHOUT_MAPPING jobs without valid company mapping${NC}"
    echo ""
    
    read -p "Do you want to fetch company data and sync? (y/n) " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Sync cancelled${NC}"
        exit 0
    fi
fi

echo ""
echo -e "${BLUE}Step 2: Fetching companies from company-service...${NC}"

# Get all companies
COMPANIES=$(curl -s -X GET "$API_GATEWAY/api/v1/companies?limit=1000" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

COMPANY_COUNT=$(echo "$COMPANIES" | jq '.data | length')
echo -e "${GREEN}✓ Found $COMPANY_COUNT companies${NC}"

if [ "$COMPANY_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}No companies found to sync${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}Step 3: Syncing company mappings...${NC}"

SUCCESS_COUNT=0
ERROR_COUNT=0

# Sync each company
echo "$COMPANIES" | jq -c '.data[]' | while read company; do
    COMPANY_ID=$(echo "$company" | jq -r '.id')
    USER_ID=$(echo "$company" | jq -r '.user_id')
    COMPANY_NAME=$(echo "$company" | jq -r '.name')
    
    echo -n "Syncing: $COMPANY_NAME (UserID: $USER_ID, CompanyID: $COMPANY_ID)... "
    
    SYNC_RESPONSE=$(curl -s -X POST "$API_GATEWAY/api/v1/admin/sync-company-mapping" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": $USER_ID,
            \"company_id\": $COMPANY_ID,
            \"company_name\": \"$COMPANY_NAME\"
        }")
    
    if echo "$SYNC_RESPONSE" | jq -e '.success' > /dev/null; then
        echo -e "${GREEN}✓${NC}"
        ((SUCCESS_COUNT++))
    else
        echo -e "${RED}✗${NC}"
        echo "  Error: $(echo "$SYNC_RESPONSE" | jq -r '.message')"
        ((ERROR_COUNT++))
    fi
done

echo ""
echo "================================"
echo -e "${GREEN}✅ Sync Complete${NC}"
echo "  Success: $SUCCESS_COUNT"
echo "  Errors: $ERROR_COUNT"
echo ""

# Final health check
echo -e "${BLUE}Step 4: Final health check...${NC}"
FINAL_CHECK=$(curl -s -X GET "$API_GATEWAY/api/v1/admin/health/data-consistency" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

echo "$FINAL_CHECK" | jq '.'

FINAL_STATUS=$(echo "$FINAL_CHECK" | jq -r '.status')
if [ "$FINAL_STATUS" == "healthy" ]; then
    echo -e "${GREEN}✓ System is now healthy!${NC}"
else
    echo -e "${YELLOW}⚠ Some issues remain, check the output above${NC}"
fi
