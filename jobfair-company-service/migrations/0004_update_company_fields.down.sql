-- Rollback changes to companies table

-- 1. Drop index for industry
DROP INDEX IF EXISTS idx_companies_industry;

-- 2. Revert company_size constraint
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_company_size_check;
ALTER TABLE companies ADD CONSTRAINT companies_company_size_check 
    CHECK (company_size IN ('1-10', '11-50', '51-200', '201-500', '501-1000', '1000+'));

-- 3. Revert industry from TEXT[] back to VARCHAR
ALTER TABLE companies ALTER COLUMN industry TYPE VARCHAR(100) USING industry[1];

-- 4. Remove contact_name column
ALTER TABLE companies DROP COLUMN IF EXISTS contact_name;
