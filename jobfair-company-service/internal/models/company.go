package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CompanySize string
type SubscriptionTier string

// Company Size sesuai dengan database constraint
const (
	CompanySize1to10     CompanySize = "1-10"
	CompanySize11to50    CompanySize = "11-50"
	CompanySize51to200   CompanySize = "51-200"
	CompanySize201to500  CompanySize = "201-500"
	CompanySize501to1000 CompanySize = "501-1000"
	CompanySize1000Plus  CompanySize = "1000+"
)

const (
	SubscriptionFree    SubscriptionTier = "free"
	SubscriptionBasic   SubscriptionTier = "basic"
	SubscriptionPremium SubscriptionTier = "premium"
	SubscriptionPro     SubscriptionTier = "pro"
)

// Company represents a company profile
type Company struct {
	ID                uint             `json:"id" gorm:"primaryKey"`
	UserID            uint             `json:"user_id" gorm:"uniqueIndex;not null"`
	Name              string           `json:"name" gorm:"not null"`
	Description       string           `json:"description" gorm:"type:text"`
	Industry          pq.StringArray   `json:"industry" gorm:"type:text[]"` // ✅ Changed to array for multiple selection
	CompanySize       CompanySize      `json:"company_size"`
	FoundedYear       int              `json:"founded_year"`
	
	// Contact Information
	ContactName string `json:"contact_name"` // ✅ Added for Contact Name from UI
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Website     string `json:"website"`
	
	// Address
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	
	// Media
	LogoURL     string         `json:"logo_url"`
	BannerURL   string         `json:"banner_url"`
	VideoURLs   pq.StringArray `json:"video_urls" gorm:"type:text[]"`
	GalleryURLs pq.StringArray `json:"gallery_urls" gorm:"type:text[]"`
	
	// Social Media
	LinkedinURL  string `json:"linkedin_url"`
	FacebookURL  string `json:"facebook_url"`
	TwitterURL   string `json:"twitter_url"`
	InstagramURL string `json:"instagram_url"`
	
	// Verification & Status
	IsVerified        bool             `json:"is_verified" gorm:"default:false"`
	VerifiedAt        *time.Time       `json:"verified_at"`
	VerificationBadge string           `json:"verification_badge"`
	IsFeatured        bool             `json:"is_featured" gorm:"default:false"`
	IsPremium         bool             `json:"is_premium" gorm:"default:false"`
	SubscriptionTier  SubscriptionTier `json:"subscription_tier" gorm:"default:'free'"`
	
	// SEO
	Slug            string         `json:"slug" gorm:"uniqueIndex"`
	MetaTitle       string         `json:"meta_title"`
	MetaDescription string         `json:"meta_description"`
	Tags            pq.StringArray `json:"tags" gorm:"type:text[]"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateCompanyRequest for creating a new company
type CreateCompanyRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Industry    []string    `json:"industry"` // ✅ Array for multiple industries
	CompanySize CompanySize `json:"company_size"`
	FoundedYear int         `json:"founded_year"`
	ContactName string      `json:"contact_name"` // ✅ Added
	Email       string      `json:"email" binding:"required,email"`
	Phone       string      `json:"phone"`
	Website     string      `json:"website"`
	Address     string      `json:"address"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
	LinkedinURL string      `json:"linkedin_url"` // ✅ Added for initial setup
}

// UpdateCompanyRequest for updating company information
type UpdateCompanyRequest struct {
	Name         *string      `json:"name"`
	Description  *string      `json:"description"`
	Industry     []string     `json:"industry"` // ✅ Array for multiple industries
	CompanySize  *CompanySize `json:"company_size"`
	FoundedYear  *int         `json:"founded_year"`
	ContactName  *string      `json:"contact_name"` // ✅ Added
	Email        *string      `json:"email"`
	Phone        *string      `json:"phone"`
	Website      *string      `json:"website"`
	Address      *string      `json:"address"`
	City         *string      `json:"city"`
	State        *string      `json:"state"`
	Country      *string      `json:"country"`
	PostalCode   *string      `json:"postal_code"`
	LinkedinURL  *string      `json:"linkedin_url"`
	FacebookURL  *string      `json:"facebook_url"`
	TwitterURL   *string      `json:"twitter_url"`
	InstagramURL *string      `json:"instagram_url"`
}

// CompanyAnalytics tracks company metrics
type CompanyAnalytics struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CompanyID    uint      `json:"company_id" gorm:"uniqueIndex;not null"`
	BoothVisits  int       `json:"booth_visits" gorm:"default:0"`
	ProfileViews int       `json:"profile_views" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CompanyMedia for company media assets
type CompanyMedia struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CompanyID    uint           `json:"company_id" gorm:"not null;index"`
	MediaType    string         `json:"media_type"`
	MediaURL     string         `json:"media_url" gorm:"not null"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	DisplayOrder int            `json:"display_order" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
