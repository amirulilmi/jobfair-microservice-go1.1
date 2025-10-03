CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    skill_name VARCHAR(255) NOT NULL,
    skill_type VARCHAR(50) NOT NULL,
    proficiency_level VARCHAR(50),
    years_of_experience INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_skills_profile_id ON skills(profile_id);
CREATE INDEX idx_skills_skill_type ON skills(skill_type);
