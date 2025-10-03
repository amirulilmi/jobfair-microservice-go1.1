package models

import (
	"time"
)

type CareerPreference struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID          uint       `gorm:"not null;uniqueIndex" json:"profile_id"`
	IsActivelyLooking  bool       `gorm:"type:boolean;default:false" json:"is_actively_looking"`
	ExpectedSalaryMin  *int       `gorm:"type:int" json:"expected_salary_min"`
	ExpectedSalaryMax  *int       `gorm:"type:int" json:"expected_salary_max"`
	SalaryCurrency     string     `gorm:"type:varchar(10);default:'IDR'" json:"salary_currency"`
	IsNegotiable       bool       `gorm:"type:boolean;default:true" json:"is_negotiable"`
	PreferredWorkTypes string     `gorm:"type:varchar(255)" json:"preferred_work_types"` // onsite,remote,hybrid (comma-separated)
	PreferredLocations string     `gorm:"type:text" json:"preferred_locations"`          // comma-separated cities
	WillingToRelocate  bool       `gorm:"type:boolean;default:false" json:"willing_to_relocate"`
	AvailableStartDate *time.Time `gorm:"type:date" json:"available_start_date"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type PositionPreference struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID    uint      `gorm:"not null;index" json:"profile_id"`
	PositionName string    `gorm:"type:varchar(255);not null" json:"position_name"`
	Priority     int       `gorm:"type:int;default:1" json:"priority"` // 1 = highest priority
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CareerPreferenceRequest struct {
	IsActivelyLooking  bool       `json:"is_actively_looking"`
	ExpectedSalaryMin  *int       `json:"expected_salary_min"`
	ExpectedSalaryMax  *int       `json:"expected_salary_max"`
	SalaryCurrency     string     `json:"salary_currency"`
	IsNegotiable       bool       `json:"is_negotiable"`
	PreferredWorkTypes []string   `json:"preferred_work_types"` // Will be joined as comma-separated
	PreferredLocations []string   `json:"preferred_locations"`  // Will be joined as comma-separated
	WillingToRelocate  bool       `json:"willing_to_relocate"`
	AvailableStartDate *time.Time `json:"available_start_date"`
}

type PositionPreferenceRequest struct {
	PositionName string `json:"position_name" binding:"required"`
	Priority     int    `json:"priority"`
}

type BulkPositionPreferenceRequest struct {
	Positions []PositionPreferenceRequest `json:"positions" binding:"required"`
}
