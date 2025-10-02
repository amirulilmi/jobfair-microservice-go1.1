-- Create job_applications table
CREATE TABLE IF NOT EXISTS job_applications (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,

    -- Application Data
    cv_url VARCHAR(500),
    cover_letter TEXT,

    -- Status
    status VARCHAR(50) DEFAULT 'applied' CHECK (status IN ('applied', 'reviewing', 'shortlisted', 'interview', 'hired', 'rejected')),
    status_note TEXT,

    -- Tracking
    viewed_at TIMESTAMP,
    reviewed_at TIMESTAMP,
    interview_at TIMESTAMP,
    responded_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT job_applications_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_job_application UNIQUE (job_id, user_id, deleted_at)
);

-- Create indexes
CREATE INDEX idx_job_applications_job_id ON job_applications(job_id);
CREATE INDEX idx_job_applications_user_id ON job_applications(user_id);
CREATE INDEX idx_job_applications_status ON job_applications(status);
CREATE INDEX idx_job_applications_created_at ON job_applications(created_at);
CREATE INDEX idx_job_applications_deleted_at ON job_applications(deleted_at);

-- Create trigger for updated_at
CREATE TRIGGER update_job_applications_updated_at
BEFORE UPDATE ON job_applications
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE job_applications IS 'Job applications from job seekers';
COMMENT ON COLUMN job_applications.status IS 'Current status of the application';
