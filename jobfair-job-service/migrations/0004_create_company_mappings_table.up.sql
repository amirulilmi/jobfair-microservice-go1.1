-- Create company_mappings table for storing user_id -> company_id mapping
CREATE TABLE IF NOT EXISTS company_mappings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    company_id INTEGER NOT NULL,
    company_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_company_mappings_user_id ON company_mappings(user_id);
CREATE INDEX idx_company_mappings_company_id ON company_mappings(company_id);

COMMENT ON TABLE company_mappings IS 'Mapping of user_id to company_id from company-service events';
