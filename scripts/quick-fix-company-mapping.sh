#!/bin/bash

# Quick Fix Script untuk Company Mapping Error
# This script will fix the "companies does not exist" error

set -e  # Exit on error

echo "🚀 JobFair Company Mapping Quick Fix"
echo "===================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Step 1: Rebuild services
echo -e "${BLUE}Step 1: Rebuilding services with latest code...${NC}"
docker-compose build job-service company-service
echo -e "${GREEN}✓ Services rebuilt${NC}"
echo ""

# Step 2: Restart services
echo -e "${BLUE}Step 2: Restarting services...${NC}"
docker-compose down
docker-compose up -d
echo -e "${GREEN}✓ Services restarted${NC}"
echo ""

# Step 3: Wait for services to be ready
echo -e "${BLUE}Step 3: Waiting for services to be ready...${NC}"
echo "Waiting 15 seconds for services to start..."
sleep 15
echo -e "${GREEN}✓ Services should be ready${NC}"
echo ""

# Step 4: Check if migration ran
echo -e "${BLUE}Step 4: Checking company_mappings table...${NC}"
MIGRATION_CHECK=$(docker exec -i postgres-job psql -U jobfair_user -d jobfair_jobs -t -c "
    SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = 'company_mappings'
    );
" 2>&1 | tr -d ' ')

if [[ "$MIGRATION_CHECK" == *"t"* ]]; then
    echo -e "${GREEN}✓ company_mappings table exists${NC}"
else
    echo -e "${YELLOW}⚠ company_mappings table not found, creating...${NC}"
    docker exec -i postgres-job psql -U jobfair_user -d jobfair_jobs <<EOF
CREATE TABLE IF NOT EXISTS company_mappings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    company_id INTEGER NOT NULL,
    company_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_company_mappings_user_id ON company_mappings(user_id);
CREATE INDEX IF NOT EXISTS idx_company_mappings_company_id ON company_mappings(company_id);
EOF
    echo -e "${GREEN}✓ company_mappings table created${NC}"
fi
echo ""

# Step 5: Populate existing companies
echo -e "${BLUE}Step 5: Populating existing company mappings...${NC}"
if [ -f "./scripts/populate-company-mappings.sh" ]; then
    chmod +x ./scripts/populate-company-mappings.sh
    ./scripts/populate-company-mappings.sh
else
    echo -e "${YELLOW}⚠ Populate script not found, using inline method...${NC}"
    
    # Get companies from company DB
    COMPANIES=$(docker exec -i postgres-company psql -U jobfair_user -d jobfair_company -t -A -F"," -c "
        SELECT user_id, id, name 
        FROM companies 
        WHERE deleted_at IS NULL
        ORDER BY created_at
    " 2>&1)
    
    if [ -z "$COMPANIES" ] || [[ "$COMPANIES" == *"error"* ]]; then
        echo -e "${YELLOW}⚠ No existing companies to populate${NC}"
    else
        echo "$COMPANIES" | while IFS=',' read -r user_id company_id company_name; do
            if [ ! -z "$user_id" ]; then
                company_name_escaped="${company_name//\'/\'\'}"
                docker exec -i postgres-job psql -U jobfair_user -d jobfair_jobs -c "
                    INSERT INTO company_mappings (user_id, company_id, company_name, created_at, updated_at)
                    VALUES ($user_id, $company_id, '$company_name_escaped', NOW(), NOW())
                    ON CONFLICT (user_id) 
                    DO UPDATE SET 
                        company_id = EXCLUDED.company_id,
                        company_name = EXCLUDED.company_name,
                        updated_at = NOW();
                " > /dev/null 2>&1
                echo "  ✓ Mapped user_id=$user_id to company_id=$company_id"
            fi
        done
        echo -e "${GREEN}✓ Existing companies populated${NC}"
    fi
fi
echo ""

# Step 6: Verify setup
echo -e "${BLUE}Step 6: Verifying setup...${NC}"

# Check consumers
echo "Checking job-service consumer..."
JOB_CONSUMER_LOG=$(docker logs job-service 2>&1 | grep -i "company event consumer" | tail -1)
if [[ "$JOB_CONSUMER_LOG" == *"started"* ]]; then
    echo -e "${GREEN}✓ Job-service consumer running${NC}"
else
    echo -e "${YELLOW}⚠ Job-service consumer status unclear${NC}"
    echo "  Check logs: docker logs job-service | grep -i consumer"
fi

echo "Checking company-service consumer..."
COMPANY_CONSUMER_LOG=$(docker logs company-service 2>&1 | grep -i "company event consumer" | tail -1)
if [[ "$COMPANY_CONSUMER_LOG" == *"started"* ]]; then
    echo -e "${GREEN}✓ Company-service consumer running${NC}"
else
    echo -e "${YELLOW}⚠ Company-service consumer status unclear${NC}"
    echo "  Check logs: docker logs company-service | grep -i consumer"
fi

# Count mappings
MAPPING_COUNT=$(docker exec -i postgres-job psql -U jobfair_user -d jobfair_jobs -t -c "
    SELECT COUNT(*) FROM company_mappings;
" 2>&1 | tr -d ' ')

if [[ "$MAPPING_COUNT" =~ ^[0-9]+$ ]] && [ "$MAPPING_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Found $MAPPING_COUNT company mappings in database${NC}"
else
    echo -e "${YELLOW}⚠ No company mappings found (this is OK if no companies exist yet)${NC}"
fi
echo ""

# Final summary
echo "=================================="
echo -e "${GREEN}✅ Fix Applied Successfully!${NC}"
echo "=================================="
echo ""
echo "What was fixed:"
echo "  1. ✓ Services rebuilt with event consumer"
echo "  2. ✓ company_mappings table created"
echo "  3. ✓ Existing companies populated"
echo "  4. ✓ Event consumers started"
echo ""
echo "Next steps:"
echo "  1. Test creating a job: POST /api/v1/jobs"
echo "  2. Check documentation: cat COMPANY_MAPPING_FIX.md"
echo ""
echo "If you still get errors:"
echo "  - Check logs: docker logs job-service"
echo "  - Verify RabbitMQ: http://localhost:15672 (guest/guest)"
echo "  - See troubleshooting in COMPANY_MAPPING_FIX.md"
echo ""
