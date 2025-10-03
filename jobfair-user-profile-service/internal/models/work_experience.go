package models

import (
	"time"
)

type WorkExperience struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID      uint       `gorm:"not null;index" json:"profile_id"`
	CompanyName    string     `gorm:"type:varchar(255);not null" json:"company_name"`
	JobPosition    string     `gorm:"type:varchar(255);not null" json:"job_position"`
	StartDate      time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate        *time.Time `gorm:"type:date" json:"end_date"`
	IsCurrentJob   bool       `gorm:"type:boolean;default:false" json:"is_current_job"`
	JobDescription string     `gorm:"type:text" json:"job_description"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type WorkExperienceRequest struct {
	CompanyName    string    `json:"company_name" binding:"required"`
	JobPosition    string    `json:"job_position" binding:"required"`
	StartDate      *DateOnly `json:"start_date" binding:"required"`
	EndDate        *DateOnly `json:"end_date"`
	IsCurrentJob   bool      `json:"is_current_job"`
	JobDescription string    `json:"job_description"`
}
