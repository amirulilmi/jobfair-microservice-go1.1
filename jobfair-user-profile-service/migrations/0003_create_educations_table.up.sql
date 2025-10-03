CREATE TABLE IF NOT EXISTS educations (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    university VARCHAR(255) NOT NULL,
    major VARCHAR(255) NOT NULL,
    degree VARCHAR(100),
    start_date DATE NOT NULL,
    end_date DATE,
    is_current BOOLEAN DEFAULT FALSE,
    gpa DECIMAL(3,2),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_educations_profile_id ON educations(profile_id);
CREATE INDEX idx_educations_is_current ON educations(is_current);
