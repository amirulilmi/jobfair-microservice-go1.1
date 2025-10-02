CREATE TABLE IF NOT EXISTS position_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL,
    position_name VARCHAR(255) NOT NULL,
    priority INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX idx_position_preferences_profile_id ON position_preferences(profile_id);
CREATE INDEX idx_position_preferences_priority ON position_preferences(priority);
