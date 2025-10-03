CREATE TABLE IF NOT EXISTS cv_documents (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL UNIQUE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    uploaded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_cv_documents_profile_id ON cv_documents(profile_id);
CREATE INDEX idx_cv_documents_uploaded_at ON cv_documents(uploaded_at);
