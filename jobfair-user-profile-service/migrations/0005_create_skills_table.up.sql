CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL,
    skill_name VARCHAR(255) NOT NULL,
    skill_type VARCHAR(50) NOT NULL,
    proficiency_level VARCHAR(50),
    years_of_experience INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_skills_profile_id ON skills(profile_id);
CREATE INDEX idx_skills_name ON skills(skill_name);
CREATE INDEX idx_skills_type ON skills(skill_type);
