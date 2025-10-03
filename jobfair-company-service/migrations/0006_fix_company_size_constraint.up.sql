-- Remove strict check constraint on company_size
ALTER TABLE companies 
DROP CONSTRAINT IF EXISTS companies_company_size_check;

-- Make company_size nullable and set default
ALTER TABLE companies 
ALTER COLUMN company_size DROP NOT NULL;

ALTER TABLE companies 
ALTER COLUMN company_size SET DEFAULT '1-10';

-- Add back a more lenient constraint (allows NULL and empty string)
ALTER TABLE companies 
ADD CONSTRAINT companies_company_size_check 
CHECK (
    company_size IS NULL OR 
    company_size = '' OR 
    company_size IN ('1-10', '11-50', '51-200', '201-500', '501-1000', '1000+')
);

COMMENT ON COLUMN companies.company_size IS 'Company size range. Nullable with default 1-10';
