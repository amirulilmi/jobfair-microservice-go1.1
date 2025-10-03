#!/bin/bash

# Script to populate existing company mappings from company-service to job-service
# This is needed for existing companies that were created before the event system was updated

echo "🔧 Populating Company Mappings to Job Service..."
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ psql command not found. Please install PostgreSQL client.${NC}"
    exit 1
fi

# Database connection strings
COMPANY_DB="postgresql://jobfair_user:jobfair_pass@localhost:5433/jobfair_company"
JOB_DB="postgresql://jobfair_user:jobfair_pass@localhost:5435/jobfair_jobs"

echo "📊 Fetching companies from company-service database..."
echo ""

# Get companies from company-service
COMPANIES=$(psql "$COMPANY_DB" -t -A -F"," -c "
    SELECT user_id, id, name 
    FROM companies 
    WHERE deleted_at IS NULL
    ORDER BY created_at
")

if [ -z "$COMPANIES" ]; then
    echo -e "${YELLOW}⚠️  No companies found in company-service database${NC}"
    exit 0
fi

# Count companies
COUNT=$(echo "$COMPANIES" | wc -l | tr -d ' ')
echo -e "${GREEN}Found $COUNT companies${NC}"
echo ""

# Insert into job-service database
echo "💾 Populating company_mappings in job-service database..."
echo ""

SUCCESS=0
FAILED=0

while IFS=',' read -r user_id company_id company_name; do
    # Skip empty lines
    if [ -z "$user_id" ]; then
        continue
    fi
    
    echo -n "  Processing company: $company_name (user_id=$user_id, company_id=$company_id)... "
    
    # Escape single quotes in company name
    company_name_escaped="${company_name//\'/\'\'}"
    
    # Insert or update mapping
    RESULT=$(psql "$JOB_DB" -t -A -c "
        INSERT INTO company_mappings (user_id, company_id, company_name, created_at, updated_at)
        VALUES ($user_id, $company_id, '$company_name_escaped', NOW(), NOW())
        ON CONFLICT (user_id) 
        DO UPDATE SET 
            company_id = EXCLUDED.company_id,
            company_name = EXCLUDED.company_name,
            updated_at = NOW()
        RETURNING id;
    " 2>&1)
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC}"
        ((SUCCESS++))
    else
        echo -e "${RED}✗${NC}"
        echo "    Error: $RESULT"
        ((FAILED++))
    fi
done <<< "$COMPANIES"

echo ""
echo "📈 Summary:"
echo -e "  ${GREEN}Success: $SUCCESS${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "  ${RED}Failed: $FAILED${NC}"
fi
echo ""

if [ $SUCCESS -gt 0 ]; then
    echo -e "${GREEN}✅ Company mappings populated successfully!${NC}"
    echo ""
    echo "You can now create jobs. The mapping table has been populated."
else
    echo -e "${RED}❌ Failed to populate company mappings${NC}"
    exit 1
fi
