-- Create jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) UNIQUE,

    -- Employment Details
    employment_type VARCHAR(50) CHECK (employment_type IN ('fulltime', 'parttime', 'contract', 'freelance', 'intern')),
    work_type VARCHAR(50) CHECK (work_type IN ('onsite', 'remote', 'hybrid')),
    experience_level VARCHAR(50) CHECK (experience_level IN ('entry', 'junior', 'mid', 'senior')),

    -- Location
    location VARCHAR(255),
    city VARCHAR(100),
    country VARCHAR(100),
    is_remote BOOLEAN DEFAULT FALSE,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Salary
    salary_min INTEGER,
    salary_max INTEGER,
    salary_currency VARCHAR(10) DEFAULT 'USD',
    salary_period VARCHAR(20) DEFAULT 'month',

    -- Requirements
    requirements TEXT[],
    responsibilities TEXT[],
    skills TEXT[],
    benefits TEXT[],

    -- Application Settings
    receive_method VARCHAR(50) DEFAULT 'email',
    contact_email VARCHAR(255),
    external_url VARCHAR(500),

    -- Metadata
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'closed', 'archived')),
    views INTEGER DEFAULT 0,
    applications INTEGER DEFAULT 0,
    deadline TIMESTAMP,
    published_at TIMESTAMP,
    closed_at TIMESTAMP,

    -- SEO
    meta_title VARCHAR(255),
    meta_description TEXT,
    tags TEXT[],

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP

    -- NOTE: No foreign key to companies table - microservices use separate databases
    -- Referential integrity is maintained at application level via API calls
);

-- Create indexes
CREATE INDEX idx_jobs_company_id ON jobs(company_id);
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_title ON jobs(title);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_employment_type ON jobs(employment_type);
CREATE INDEX idx_jobs_work_type ON jobs(work_type);
CREATE INDEX idx_jobs_experience_level ON jobs(experience_level);
CREATE INDEX idx_jobs_location ON jobs(location);
CREATE INDEX idx_jobs_city ON jobs(city);
CREATE INDEX idx_jobs_deleted_at ON jobs(deleted_at);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);

-- GIN indexes for array columns
CREATE INDEX idx_jobs_skills ON jobs USING GIN (skills);
CREATE INDEX idx_jobs_tags ON jobs USING GIN (tags);

-- Full-text search index
CREATE INDEX idx_jobs_title_fulltext ON jobs USING GIN (to_tsvector('english', title));
CREATE INDEX idx_jobs_description_fulltext ON jobs USING GIN (to_tsvector('english', description));

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_jobs_updated_at
BEFORE UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE jobs IS 'Job postings from companies';
COMMENT ON COLUMN jobs.company_id IS 'Reference to company (validated via API, not FK)';
COMMENT ON COLUMN jobs.user_id IS 'Company user who posted the job';
COMMENT ON COLUMN jobs.receive_method IS 'How to receive applications: email or external';
