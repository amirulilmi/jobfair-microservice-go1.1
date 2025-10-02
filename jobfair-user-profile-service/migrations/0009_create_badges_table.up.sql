CREATE TABLE IF NOT EXISTS badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    profile_id UUID NOT NULL,
    badge_id UUID NOT NULL,
    earned_at TIMESTAMP NOT NULL,
    PRIMARY KEY (profile_id, badge_id),
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE
);

CREATE INDEX idx_badges_type ON badges(badge_type);
CREATE INDEX idx_badges_rarity ON badges(rarity);
CREATE INDEX idx_profile_badges_profile_id ON profile_badges(profile_id);
CREATE INDEX idx_profile_badges_badge_id ON profile_badges(badge_id);

-- Insert default badges
INSERT INTO badges (badge_name, description, badge_type, points, rarity) VALUES
('Profile Complete', 'Completed 100% of profile', 'profile_completion', 100, 'common'),
('Early Adopter', 'One of the first users', 'early_adopter', 50, 'rare'),
('Skill Master', 'Added 10+ verified skills', 'skill_verified', 75, 'epic'),
('Experience Pro', 'Added 3+ work experiences', 'experience', 50, 'common'),
('Education Elite', 'Added education background', 'education', 25, 'common');
