package models

// APIResponse is the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

// CreateJobRequest is the request for creating a job
type CreateJobRequest struct {
	Title            string           `json:"title" binding:"required"`
	Description      string           `json:"description" binding:"required"`
	EmploymentType   EmploymentType   `json:"employment_type" binding:"required"`
	WorkType         WorkType         `json:"work_type" binding:"required"`
	ExperienceLevel  ExperienceLevel  `json:"experience_level" binding:"required"`
	Location         string           `json:"location" binding:"required"`
	SalaryMin        int              `json:"salary_min"`
	SalaryMax        int              `json:"salary_max"`
	Requirements     []string         `json:"requirements"`
	Responsibilities []string         `json:"responsibilities"`
	Skills           []string         `json:"skills"`
	Benefits         []string         `json:"benefits"`
	ReceiveMethod    string           `json:"receive_method"` // email or external
	ContactEmail     string           `json:"contact_email"`
	ExternalURL      *string          `json:"external_url"`
	Deadline         *string          `json:"deadline"` // ISO 8601 format
	Tags             []string         `json:"tags"`
}

// UpdateJobRequest is the request for updating a job
type UpdateJobRequest struct {
	Title            *string          `json:"title"`
	Description      *string          `json:"description"`
	EmploymentType   *EmploymentType  `json:"employment_type"`
	WorkType         *WorkType        `json:"work_type"`
	ExperienceLevel  *ExperienceLevel `json:"experience_level"`
	Location         *string          `json:"location"`
	SalaryMin        *int             `json:"salary_min"`
	SalaryMax        *int             `json:"salary_max"`
	Requirements     []string         `json:"requirements"`
	Responsibilities []string         `json:"responsibilities"`
	Skills           []string         `json:"skills"`
	Benefits         []string         `json:"benefits"`
	ReceiveMethod    *string          `json:"receive_method"`
	ContactEmail     *string          `json:"contact_email"`
	ExternalURL      *string          `json:"external_url"`
	Deadline         *string          `json:"deadline"`
	Tags             []string         `json:"tags"`
}

// ApplyJobRequest is the request for applying to a job
type ApplyJobRequest struct {
	CVURL       string `json:"cv_url"`
	CoverLetter string `json:"cover_letter"`
}

// BulkApplyRequest is the request for applying to multiple jobs
type BulkApplyRequest struct {
	JobIDs      []uint `json:"job_ids" binding:"required"`
	CVURL       string `json:"cv_url"`
	CoverLetter string `json:"cover_letter"`
}

// UpdateApplicationStatusRequest is the request for updating application status
type UpdateApplicationStatusRequest struct {
	Status     ApplicationStatus `json:"status" binding:"required"`
	StatusNote string            `json:"status_note"`
}

// JobListFilter represents filters for job listing
type JobListFilter struct {
	Search          string            `form:"search"`
	EmploymentType  []EmploymentType  `form:"employment_type"`
	WorkType        []WorkType        `form:"work_type"`
	ExperienceLevel []ExperienceLevel `form:"experience_level"`
	Location        string            `form:"location"`
	SalaryMin       int               `form:"salary_min"`
	SalaryMax       int               `form:"salary_max"`
	CompanyID       uint              `form:"company_id"`
	Tags            []string          `form:"tags"`
	Status          JobStatus         `form:"status"`
	Page            int               `form:"page"`
	Limit           int               `form:"limit"`
	OrderBy         string            `form:"order_by"` // created_at, views, applications
	Order           string            `form:"order"`    // asc, desc
}

// JobDetailResponse includes company information
type JobDetailResponse struct {
	Job         *Job                   `json:"job"`
	Company     map[string]interface{} `json:"company,omitempty"`
	IsSaved     bool                   `json:"is_saved"`
	HasApplied  bool                   `json:"has_applied"`
	Application *JobApplication        `json:"application,omitempty"`
}

// ApplicationListResponse includes job and company information
type ApplicationListResponse struct {
	Applications []ApplicationWithJob `json:"applications"`
	Meta         PaginationMeta       `json:"meta"`
}

type ApplicationWithJob struct {
	Application *JobApplication        `json:"application"`
	Job         *Job                   `json:"job"`
	Company     map[string]interface{} `json:"company"`
}

// JobWithCompany combines job and company data for list responses
type JobWithCompany struct {
	*Job
	Company map[string]interface{} `json:"company,omitempty"`
}

// JobListResponse for list jobs with company data
type JobListResponse struct {
	Jobs []JobWithCompany `json:"jobs"`
	Meta PaginationMeta  `json:"meta"`
}
