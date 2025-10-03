-- Revert contact_name and industry changes
ALTER TABLE company_basic_profiles 
DROP COLUMN IF EXISTS contact_name;

-- Revert industry back to varchar
ALTER TABLE company_basic_profiles 
ALTER COLUMN industry TYPE VARCHAR(100) USING industry[1];
