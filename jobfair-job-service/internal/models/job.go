package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Enums
type JobStatus string
type EmploymentType string
type WorkType string
type ExperienceLevel string
type ApplicationStatus string

const (
	JobStatusDraft     JobStatus = "draft"
	JobStatusPublished JobStatus = "published"
	JobStatusClosed    JobStatus = "closed"
	JobStatusArchived  JobStatus = "archived"
)

const (
	EmploymentTypeFullTime  EmploymentType = "fulltime"
	EmploymentTypePartTime  EmploymentType = "parttime"
	EmploymentTypeContract  EmploymentType = "contract"
	EmploymentTypeFreelance EmploymentType = "freelance"
	EmploymentTypeIntern    EmploymentType = "intern"
)

const (
	WorkTypeOnsite WorkType = "onsite"
	WorkTypeRemote WorkType = "remote"
	WorkTypeHybrid WorkType = "hybrid"
)

const (
	ExperienceLevelEntry  ExperienceLevel = "entry"
	ExperienceLevelJunior ExperienceLevel = "junior"
	ExperienceLevelMid    ExperienceLevel = "mid"
	ExperienceLevelSenior ExperienceLevel = "senior"
)

const (
	ApplicationStatusApplied     ApplicationStatus = "applied"
	ApplicationStatusReviewing   ApplicationStatus = "reviewing"
	ApplicationStatusShortlisted ApplicationStatus = "shortlisted"
	ApplicationStatusInterview   ApplicationStatus = "interview"
	ApplicationStatusHired       ApplicationStatus = "hired"
	ApplicationStatusRejected    ApplicationStatus = "rejected"
)

// Job represents a job posting
type Job struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	CompanyID   uint   `json:"company_id" gorm:"not null;index"`
	UserID      uint   `json:"user_id" gorm:"not null;index"` // Company user who posted
	Title       string `json:"title" gorm:"not null;index"`
	Description string `json:"description" gorm:"type:text"`
	Slug        string `json:"slug" gorm:"uniqueIndex"`

	// Employment Details
	EmploymentType  EmploymentType  `json:"employment_type" gorm:"type:varchar(50)"`
	WorkType        WorkType        `json:"work_type" gorm:"type:varchar(50)"`
	ExperienceLevel ExperienceLevel `json:"experience_level" gorm:"type:varchar(50)"`

	// Location
	Location    string  `json:"location"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	IsRemote    bool    `json:"is_remote" gorm:"default:false"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`

	// Salary
	SalaryMin      int    `json:"salary_min"`
	SalaryMax      int    `json:"salary_max"`
	SalaryCurrency string `json:"salary_currency" gorm:"default:'USD'"`
	SalaryPeriod   string `json:"salary_period" gorm:"default:'month'"` // month, year

	// Requirements
	Requirements     pq.StringArray `json:"requirements" gorm:"type:text[]"`
	Responsibilities pq.StringArray `json:"responsibilities" gorm:"type:text[]"`
	Skills           pq.StringArray `json:"skills" gorm:"type:text[]"`
	Benefits         pq.StringArray `json:"benefits" gorm:"type:text[]"`

	// Application Settings
	ReceiveMethod string  `json:"receive_method" gorm:"default:'email'"` // email, external
	ContactEmail  string  `json:"contact_email"`
	ExternalURL   *string `json:"external_url"`

	// Metadata
	Status        JobStatus  `json:"status" gorm:"default:'draft'"`
	Views         int        `json:"views" gorm:"default:0"`
	Applications  int        `json:"applications" gorm:"default:0"`
	Deadline      *time.Time `json:"deadline"`
	PublishedAt   *time.Time `json:"published_at"`
	ClosedAt      *time.Time `json:"closed_at"`

	// SEO
	MetaTitle       string         `json:"meta_title"`
	MetaDescription string         `json:"meta_description"`
	Tags            pq.StringArray `json:"tags" gorm:"type:text[]"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// JobApplication represents a job application
type JobApplication struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	JobID  uint `json:"job_id" gorm:"not null;index"`
	UserID uint `json:"user_id" gorm:"not null;index"` // Job seeker

	// Application Data
	CVURL       string `json:"cv_url"`
	CoverLetter string `json:"cover_letter" gorm:"type:text"`

	// Status
	Status     ApplicationStatus `json:"status" gorm:"default:'applied'"`
	StatusNote string            `json:"status_note" gorm:"type:text"` // Notes from recruiter

	// Tracking
	ViewedAt    *time.Time `json:"viewed_at"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	InterviewAt *time.Time `json:"interview_at"`
	RespondedAt *time.Time `json:"responded_at"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations (not stored in DB)
	Job *Job `json:"job,omitempty" gorm:"foreignKey:JobID"`
}

// SavedJob represents a bookmarked job
type SavedJob struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	JobID     uint           `json:"job_id" gorm:"not null;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Job *Job `json:"job,omitempty" gorm:"foreignKey:JobID"`
}

// Table names
func (Job) TableName() string {
	return "jobs"
}

func (JobApplication) TableName() string {
	return "job_applications"
}

func (SavedJob) TableName() string {
	return "saved_jobs"
}
