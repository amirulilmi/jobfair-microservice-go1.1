-- Add contact_name and change industry to array
ALTER TABLE company_basic_profiles 
ADD COLUMN IF NOT EXISTS contact_name VARCHAR(255);

-- Change industry to text array for multiple industries
ALTER TABLE company_basic_profiles 
ALTER COLUMN industry TYPE TEXT[] USING ARRAY[industry]::TEXT[];

COMMENT ON COLUMN company_basic_profiles.contact_name IS 'Contact person name at the company';
COMMENT ON COLUMN company_basic_profiles.industry IS 'Multiple industries as array';
