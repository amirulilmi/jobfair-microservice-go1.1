package models

import (
	"time"
)

type Badge struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BadgeName   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"badge_name"`
	Description string    `gorm:"type:text" json:"description"`
	IconURL     string    `gorm:"type:varchar(500)" json:"icon_url"`
	BadgeType   string    `gorm:"type:varchar(50);not null" json:"badge_type"` // profile_completion, skill_verified, early_adopter, etc.
	Points      int       `gorm:"type:int;default:0" json:"points"`
	Rarity      string    `gorm:"type:varchar(50)" json:"rarity"` // common, rare, epic, legendary
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ProfileBadge struct {
	ProfileID uint      `gorm:"primaryKey" json:"profile_id"`
	BadgeID   uint      `gorm:"primaryKey" json:"badge_id"`
	EarnedAt  time.Time `gorm:"type:timestamp;not null" json:"earned_at"`
	Badge     Badge     `gorm:"foreignKey:BadgeID" json:"badge"`
}
