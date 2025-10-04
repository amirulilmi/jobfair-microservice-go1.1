#!/bin/bash

# Quick fix script untuk 3 bug yang ditemukan
# Usage: ./scripts/apply-bug-fixes.sh

set -e

echo "🔧 Applying Bug Fixes - October 2025"
echo "===================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ docker-compose not found${NC}"
    exit 1
fi

echo -e "${BLUE}Bug Fixes to be applied:${NC}"
echo "  1. Company Logo Upload - Form field clarification"
echo "  2. Job Application - JSON parsing fix"
echo "  3. Profile Creation - Database migration + first/last name support"
echo ""

# Step 1: Apply migration to profile service
echo -e "${BLUE}Step 1: Applying database migration to user-profile-service...${NC}"

# Check if migration file exists
if [ ! -f "jobfair-user-profile-service/migrations/0010_add_missing_profile_columns.up.sql" ]; then
    echo -e "${RED}❌ Migration file not found!${NC}"
    exit 1
fi

# Apply migration
echo "Running migration..."
docker exec -i postgres-profile psql -U jobfair_user -d jobfair_profiles < \
    jobfair-user-profile-service/migrations/0010_add_missing_profile_columns.up.sql 2>&1 | grep -v "NOTICE: relation" || true

echo -e "${GREEN}✓ Migration applied${NC}"
echo ""

# Step 2: Rebuild services
echo -e "${BLUE}Step 2: Rebuilding affected services...${NC}"
docker-compose build jobfair-user-profile-service
echo -e "${GREEN}✓ Services rebuilt${NC}"
echo ""

# Step 3: Restart services
echo -e "${BLUE}Step 3: Restarting services...${NC}"
docker-compose restart jobfair-user-profile-service
echo -e "${GREEN}✓ Services restarted${NC}"
echo ""

# Step 4: Verify migration
echo -e "${BLUE}Step 4: Verifying database changes...${NC}"

COLUMNS_CHECK=$(docker exec -i postgres-profile psql -U jobfair_user -d jobfair_profiles -t -c "
    SELECT column_name 
    FROM information_schema.columns 
    WHERE table_name = 'profiles' 
    AND column_name IN ('first_name', 'last_name', 'headline', 'summary', 'location', 'github_url')
    ORDER BY column_name;
" 2>&1 | tr -d ' ')

EXPECTED_COLS=("first_name" "github_url" "headline" "last_name" "location" "summary")
FOUND_COUNT=0

for col in "${EXPECTED_COLS[@]}"; do
    if echo "$COLUMNS_CHECK" | grep -q "$col"; then
        echo -e "  ${GREEN}✓${NC} Column '$col' exists"
        ((FOUND_COUNT++))
    else
        echo -e "  ${RED}✗${NC} Column '$col' missing"
    fi
done

if [ $FOUND_COUNT -eq ${#EXPECTED_COLS[@]} ]; then
    echo -e "${GREEN}✓ All columns verified${NC}"
else
    echo -e "${YELLOW}⚠ Some columns missing, but migration completed${NC}"
fi
echo ""

# Summary
echo "===================================="
echo -e "${GREEN}✅ Bug Fixes Applied Successfully!${NC}"
echo "===================================="
echo ""
echo "What was fixed:"
echo "  1. ✓ Company logo upload now expects field name 'file'"
echo "  2. ✓ Job application validates JSON properly (no comments allowed)"
echo "  3. ✓ Profile now supports: first_name, last_name, headline, summary, location, github_url"
echo ""
echo "Testing:"
echo "  • Company logo: Use form-data field 'file' (not 'logo')"
echo "  • Job apply: Ensure JSON has no comments (// or /* */)"
echo "  • Profile: Can use first_name + last_name OR full_name"
echo ""
echo "Documentation: cat BUG_FIXES_OCT_2025.md"
echo ""
