package models

import (
	"time"

	"github.com/google/uuid"
)

type Education struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProfileID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"profile_id"`
	University  string     `gorm:"type:varchar(255);not null" json:"university"`
	Major       string     `gorm:"type:varchar(255);not null" json:"major"`
	Degree      string     `gorm:"type:varchar(100)" json:"degree"` // Bachelor, Master, PhD, etc.
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date"`
	IsCurrent   bool       `gorm:"type:boolean;default:false" json:"is_current"`
	GPA         *float64   `gorm:"type:decimal(3,2)" json:"gpa"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type EducationRequest struct {
	University  string     `json:"university" binding:"required"`
	Major       string     `json:"major" binding:"required"`
	Degree      string     `json:"degree"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	IsCurrent   bool       `json:"is_current"`
	GPA         *float64   `json:"gpa"`
	Description string     `json:"description"`
}
