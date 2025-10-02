CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    full_name VARCHAR(255),
    phone_number VARCHAR(20),
    bio TEXT,
    profile_picture_url VARCHAR(500),
    banner_image_url VARCHAR(500),
    date_of_birth DATE,
    gender VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Indonesia',
    postal_code VARCHAR(10),
    linkedin_url VARCHAR(255),
    portfolio_url VARCHAR(255),
    completion_status INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_city ON profiles(city);
CREATE INDEX idx_profiles_completion_status ON profiles(completion_status);
