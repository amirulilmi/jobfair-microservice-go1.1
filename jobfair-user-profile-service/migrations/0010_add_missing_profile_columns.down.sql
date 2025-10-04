-- Remove added columns from profiles table
ALTER TABLE profiles 
DROP COLUMN IF EXISTS github_url,
DROP COLUMN IF EXISTS location,
DROP COLUMN IF EXISTS summary,
DROP COLUMN IF EXISTS headline,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS first_name;

-- Drop indexes
DROP INDEX IF EXISTS idx_profiles_location;
DROP INDEX IF EXISTS idx_profiles_last_name;
DROP INDEX IF EXISTS idx_profiles_first_name;
