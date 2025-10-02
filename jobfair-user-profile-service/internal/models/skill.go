package models

import (
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProfileID     uuid.UUID `gorm:"type:uuid;not null;index" json:"profile_id"`
	SkillName     string    `gorm:"type:varchar(255);not null" json:"skill_name"`
	SkillType     string    `gorm:"type:varchar(50);not null" json:"skill_type"` // technical, soft
	ProficiencyLevel string `gorm:"type:varchar(50)" json:"proficiency_level"` // beginner, intermediate, advanced, expert
	YearsOfExperience *int   `gorm:"type:int" json:"years_of_experience"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type SkillRequest struct {
	SkillName         string `json:"skill_name" binding:"required"`
	SkillType         string `json:"skill_type" binding:"required,oneof=technical soft"`
	ProficiencyLevel  string `json:"proficiency_level" binding:"omitempty,oneof=beginner intermediate advanced expert"`
	YearsOfExperience *int   `json:"years_of_experience"`
}

type BulkSkillRequest struct {
	TechnicalSkills []SkillRequest `json:"technical_skills"`
	SoftSkills      []SkillRequest `json:"soft_skills"`
}
