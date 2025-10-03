-- Revert contact_name column
ALTER TABLE companies 
DROP COLUMN IF EXISTS contact_name;

-- Revert industry back to varchar (first element only)
ALTER TABLE companies 
ALTER COLUMN industry TYPE VARCHAR(100) USING industry[1];
