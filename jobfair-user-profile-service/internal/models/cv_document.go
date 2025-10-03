package models

import (
	"time"
)

type CVDocument struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProfileID  uint      `gorm:"not null;uniqueIndex" json:"profile_id"`
	FileName   string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileURL    string    `gorm:"type:varchar(500);not null" json:"file_url"`
	FileSize   int64     `gorm:"type:bigint;not null" json:"file_size"`          // in bytes
	FileType   string    `gorm:"type:varchar(50);not null" json:"file_type"`     // pdf, doc, docx
	IsVerified bool      `gorm:"type:boolean;default:false" json:"is_verified"`
	UploadedAt time.Time `gorm:"type:timestamp;not null" json:"uploaded_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CVUploadResponse struct {
	ID         uint      `json:"id"`
	FileName   string    `json:"file_name"`
	FileURL    string    `json:"file_url"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	IsVerified bool      `json:"is_verified"`
	UploadedAt time.Time `json:"uploaded_at"`
}
