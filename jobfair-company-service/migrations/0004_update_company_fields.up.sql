-- Update companies table to match UI requirements

-- 1. Add contact_name column
ALTER TABLE companies ADD COLUMN IF NOT EXISTS contact_name VARCHAR(255);

-- 2. Change industry from VARCHAR to TEXT[] for multiple selection
ALTER TABLE companies ALTER COLUMN industry TYPE TEXT[] USING ARRAY[industry]::TEXT[];

-- 3. Update company_size constraint to match UI values
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_company_size_check;
ALTER TABLE companies ADD CONSTRAINT companies_company_size_check 
    CHECK (company_size IN ('small', 'medium', 'large'));

-- 4. Create index for industry array
CREATE INDEX IF NOT EXISTS idx_companies_industry ON companies USING GIN (industry);

-- Comments
COMMENT ON COLUMN companies.contact_name IS 'Contact person name for the company';
COMMENT ON COLUMN companies.industry IS 'Industry categories (supports multiple selection)';
COMMENT ON COLUMN companies.company_size IS 'Company size: small (1-50), medium (51-250), large (251+)';
