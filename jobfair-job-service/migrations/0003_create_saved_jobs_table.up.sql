-- Create saved_jobs table (bookmarks)
CREATE TABLE IF NOT EXISTS saved_jobs (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT saved_jobs_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_saved_job UNIQUE (job_id, user_id, deleted_at)
);

-- Create indexes
CREATE INDEX idx_saved_jobs_job_id ON saved_jobs(job_id);
CREATE INDEX idx_saved_jobs_user_id ON saved_jobs(user_id);
CREATE INDEX idx_saved_jobs_deleted_at ON saved_jobs(deleted_at);

-- Comments
COMMENT ON TABLE saved_jobs IS 'Bookmarked jobs by users';
