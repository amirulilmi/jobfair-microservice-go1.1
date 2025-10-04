-- Add missing columns to profiles table
ALTER TABLE profiles 
ADD COLUMN IF NOT EXISTS first_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS last_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS headline VARCHAR(255),
ADD COLUMN IF NOT EXISTS summary TEXT,
ADD COLUMN IF NOT EXISTS location VARCHAR(255),
ADD COLUMN IF NOT EXISTS github_url VARCHAR(255);

-- Create index for searching by name
CREATE INDEX IF NOT EXISTS idx_profiles_first_name ON profiles(first_name);
CREATE INDEX IF NOT EXISTS idx_profiles_last_name ON profiles(last_name);
CREATE INDEX IF NOT EXISTS idx_profiles_location ON profiles(location);

COMMENT ON COLUMN profiles.first_name IS 'First name of the user';
COMMENT ON COLUMN profiles.last_name IS 'Last name of the user';
COMMENT ON COLUMN profiles.headline IS 'Professional headline (e.g., "Senior Software Engineer")';
COMMENT ON COLUMN profiles.summary IS 'Professional summary or bio';
COMMENT ON COLUMN profiles.location IS 'Current location';
COMMENT ON COLUMN profiles.github_url IS 'GitHub profile URL';
