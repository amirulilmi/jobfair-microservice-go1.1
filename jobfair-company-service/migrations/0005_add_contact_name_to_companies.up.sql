-- Add contact_name column if not exists
ALTER TABLE companies 
ADD COLUMN IF NOT EXISTS contact_name VARCHAR(255);

-- Change industry to text array if not already
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'companies' 
        AND column_name = 'industry' 
        AND data_type != 'ARRAY'
    ) THEN
        ALTER TABLE companies 
        ALTER COLUMN industry TYPE TEXT[] USING ARRAY[industry]::TEXT[];
    END IF;
END $$;

COMMENT ON COLUMN companies.contact_name IS 'Contact person name at the company';
COMMENT ON COLUMN companies.industry IS 'Multiple industries as array';
