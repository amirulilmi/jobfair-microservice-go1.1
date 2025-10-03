package models

import (
	"time"
)

type Profile struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	FullName          string     `gorm:"type:varchar(255)" json:"full_name"`
	PhoneNumber       string     `gorm:"type:varchar(20)" json:"phone_number"`
	Bio               string     `gorm:"type:text" json:"bio"`
	ProfilePictureURL string     `gorm:"type:varchar(500)" json:"profile_picture_url"`
	BannerImageURL    string     `gorm:"type:varchar(500)" json:"banner_image_url"`
	DateOfBirth       *time.Time `gorm:"type:date" json:"date_of_birth"`
	Gender            string     `gorm:"type:varchar(20)" json:"gender"` // male, female, other
	Address           string     `gorm:"type:text" json:"address"`
	City              string     `gorm:"type:varchar(100)" json:"city"`
	Province          string     `gorm:"type:varchar(100)" json:"province"`
	Country           string     `gorm:"type:varchar(100);default:'Indonesia'" json:"country"`
	PostalCode        string     `gorm:"type:varchar(10)" json:"postal_code"`
	LinkedInURL       string     `gorm:"type:varchar(255)" json:"linkedin_url"`
	PortfolioURL      string     `gorm:"type:varchar(255)" json:"portfolio_url"`
	CompletionStatus  int        `gorm:"type:int;default:0" json:"completion_status"` // 0-100%
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	WorkExperiences     []WorkExperience     `gorm:"foreignKey:ProfileID" json:"work_experiences,omitempty"`
	Educations          []Education          `gorm:"foreignKey:ProfileID" json:"educations,omitempty"`
	Certifications      []Certification      `gorm:"foreignKey:ProfileID" json:"certifications,omitempty"`
	Skills              []Skill              `gorm:"foreignKey:ProfileID" json:"skills,omitempty"`
	PositionPreferences []PositionPreference `gorm:"foreignKey:ProfileID" json:"position_preferences,omitempty"`
	CareerPreference    *CareerPreference    `gorm:"foreignKey:ProfileID" json:"career_preference,omitempty"`
	CVDocument          *CVDocument          `gorm:"foreignKey:ProfileID" json:"cv_document,omitempty"`
	Badges              []Badge              `gorm:"many2many:profile_badges;" json:"badges,omitempty"`
}

type ProfileUpdateRequest struct {
	FullName          *string    `json:"full_name"`
	PhoneNumber       *string    `json:"phone_number"`
	Bio               *string    `json:"bio"`
	DateOfBirth       *time.Time `json:"date_of_birth"`
	Gender            *string    `json:"gender"`
	Address           *string    `json:"address"`
	City              *string    `json:"city"`
	Province          *string    `json:"province"`
	Country           *string    `json:"country"`
	PostalCode        *string    `json:"postal_code"`
	LinkedInURL       *string    `json:"linkedin_url"`
	PortfolioURL      *string    `json:"portfolio_url"`
	ProfilePictureURL *string    `json:"profile_picture_url"`
	BannerImageURL    *string    `json:"banner_image_url"`
}
