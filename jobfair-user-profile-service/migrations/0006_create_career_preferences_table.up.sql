CREATE TABLE IF NOT EXISTS career_preferences (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL UNIQUE,
    is_actively_looking BOOLEAN DEFAULT FALSE,
    expected_salary_min INT,
    expected_salary_max INT,
    salary_currency VARCHAR(10) DEFAULT 'IDR',
    is_negotiable BOOLEAN DEFAULT TRUE,
    preferred_work_types VARCHAR(255),
    preferred_locations TEXT,
    willing_to_relocate BOOLEAN DEFAULT FALSE,
    available_start_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_career_preferences_profile_id ON career_preferences(profile_id);
CREATE INDEX idx_career_preferences_actively_looking ON career_preferences(is_actively_looking);
