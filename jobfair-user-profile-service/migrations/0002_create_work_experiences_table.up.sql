CREATE TABLE IF NOT EXISTS work_experiences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    job_position VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_current_job BOOLEAN DEFAULT FALSE,
    job_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_work_experiences_profile_id ON work_experiences(profile_id);
CREATE INDEX idx_work_experiences_is_current ON work_experiences(is_current_job);
