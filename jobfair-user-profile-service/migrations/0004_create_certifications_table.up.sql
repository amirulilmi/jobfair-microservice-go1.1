CREATE TABLE IF NOT EXISTS certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL,
    certification_name VARCHAR(255) NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    issue_date DATE NOT NULL,
    expiry_date DATE,
    credential_id VARCHAR(255),
    credential_url VARCHAR(500),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_certifications_profile_id ON certifications(profile_id);
CREATE INDEX idx_certifications_issue_date ON certifications(issue_date);
