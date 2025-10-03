package models

import (
	"time"
)

type CareerPreference struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID          uint      `gorm:"not null;uniqueIndex" json:"profile_id"`
	IsActivelyLooking  bool      `gorm:"type:boolean;default:false" json:"is_actively_looking"`
	ExpectedSalaryMin  *int      `gorm:"type:int" json:"expected_salary_min"`
	ExpectedSalaryMax  *int      `gorm:"type:int" json:"expected_salary_max"`
	SalaryCurrency     string    `gorm:"type:varchar(10);default:'IDR'" json:"salary_currency"`
	IsNegotiable       bool      `gorm:"type:boolean;default:true" json:"is_negotiable"`
	PreferredWorkTypes string    `gorm:"type:varchar(255)" json:"preferred_work_types"` // onsite,remote,hybrid (comma-separated)
	PreferredLocations string    `gorm:"type:text" json:"preferred_locations"`          // comma-separated cities
	WillingToRelocate  bool      `gorm:"type:boolean;default:false" json:"willing_to_relocate"`
	AvailableStartDate time.Time `gorm:"type:date" json:"available_start_date"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type PositionPreference struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID    uint      `gorm:"not null;index" json:"profile_id"`
	PositionName string    `gorm:"type:varchar(255);not null" json:"position_name"`
	Priority     int       `gorm:"type:int;default:1" json:"priority"` // 1 = highest priority
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Updated CareerPreferenceRequest to match actual API request format
type CareerPreferenceRequest struct {
	JobType           string    `json:"job_type"`      // Maps to PreferredWorkTypes
	WorkLocation      string    `json:"work_location"` // Maps to PreferredWorkTypes
	ExpectedSalaryMin *int      `json:"expected_salary_min"`
	ExpectedSalaryMax *int      `json:"expected_salary_max"`
	Currency          string    `json:"currency"` // Maps to SalaryCurrency
	WillingToRelocate bool      `json:"willing_to_relocate"`
	AvailableFrom     *DateOnly `json:"available_from"` // Maps to AvailableStartDate
}

type PositionPreferenceRequest struct {
	PositionName string `json:"position_name" binding:"required"`
	Priority     int    `json:"priority"`
}

// Updated BulkPositionPreferenceRequest to support both string array and object array
type BulkPositionPreferenceRequest struct {
	Positions []string `json:"positions" binding:"required"`
}
