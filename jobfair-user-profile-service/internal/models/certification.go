package models

import (
	"time"
)

type Certification struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID         uint       `gorm:"not null;index" json:"profile_id"`
	CertificationName string     `gorm:"type:varchar(255);not null" json:"certification_name"`
	Organizer         string     `gorm:"type:varchar(255);not null" json:"organizer"`
	IssueDate         time.Time  `gorm:"type:date;not null" json:"issue_date"`
	ExpiryDate        *time.Time `gorm:"type:date" json:"expiry_date"`
	CredentialID      string     `gorm:"type:varchar(255)" json:"credential_id"`
	CredentialURL     string     `gorm:"type:varchar(500)" json:"credential_url"`
	Description       string     `gorm:"type:text" json:"description"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type CertificationRequest struct {
	CertificationName string     `json:"certification_name" binding:"required"`
	Organizer         string     `json:"organizer" binding:"required"`
	IssueDate         time.Time  `json:"issue_date" binding:"required"`
	ExpiryDate        *time.Time `json:"expiry_date"`
	CredentialID      string     `json:"credential_id"`
	CredentialURL     string     `json:"credential_url"`
	Description       string     `json:"description"`
}
