-- Revert company_size constraint changes
ALTER TABLE companies 
DROP CONSTRAINT IF EXISTS companies_company_size_check;

-- Restore strict constraint
ALTER TABLE companies 
ADD CONSTRAINT companies_company_size_check 
CHECK (company_size IN ('1-10', '11-50', '51-200', '201-500', '501-1000', '1000+'));

-- Remove default
ALTER TABLE companies 
ALTER COLUMN company_size DROP DEFAULT;
