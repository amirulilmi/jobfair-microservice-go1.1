CREATE TABLE IF NOT EXISTS badges (
    id SERIAL PRIMARY KEY,
    badge_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    icon_url VARCHAR(500),
    badge_type VARCHAR(50) NOT NULL,
    points INT DEFAULT 0,
    rarity VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profile_badges (
    profile_id INTEGER NOT NULL,
    badge_id INTEGER NOT NULL,
    earned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (profile_id, badge_id),
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE
);

CREATE INDEX idx_badges_badge_type ON badges(badge_type);
CREATE INDEX idx_badges_rarity ON badges(rarity);
CREATE INDEX idx_profile_badges_profile_id ON profile_badges(profile_id);
CREATE INDEX idx_profile_badges_badge_id ON profile_badges(badge_id);
