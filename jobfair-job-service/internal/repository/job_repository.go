package repository

import (
	"fmt"
	"jobfair-job-service/internal/models"

	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create creates a new job
func (r *JobRepository) Create(job *models.Job) (*models.Job, error) {
	if err := r.db.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

// GetByID retrieves a job by ID
func (r *JobRepository) GetByID(id uint) (*models.Job, error) {
	var job models.Job
	if err := r.db.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// GetBySlug retrieves a job by slug
func (r *JobRepository) GetBySlug(slug string) (*models.Job, error) {
	var job models.Job
	if err := r.db.Where("slug = ?", slug).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// Update updates a job
func (r *JobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

// Delete soft deletes a job
func (r *JobRepository) Delete(id uint) error {
	return r.db.Delete(&models.Job{}, id).Error
}

// List retrieves jobs with filters and pagination
func (r *JobRepository) List(filter models.JobListFilter) ([]*models.Job, int64, error) {
	var jobs []*models.Job
	var total int64

	query := r.db.Model(&models.Job{})

	// Apply filters
	if filter.Search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if len(filter.EmploymentType) > 0 {
		query = query.Where("employment_type IN ?", filter.EmploymentType)
	}

	if len(filter.WorkType) > 0 {
		query = query.Where("work_type IN ?", filter.WorkType)
	}

	if len(filter.ExperienceLevel) > 0 {
		query = query.Where("experience_level IN ?", filter.ExperienceLevel)
	}

	if filter.Location != "" {
		query = query.Where("location ILIKE ? OR city ILIKE ?", "%"+filter.Location+"%", "%"+filter.Location+"%")
	}

	if filter.SalaryMin > 0 {
		query = query.Where("salary_max >= ?", filter.SalaryMin)
	}

	if filter.SalaryMax > 0 {
		query = query.Where("salary_min <= ?", filter.SalaryMax)
	}

	if filter.CompanyID > 0 {
		query = query.Where("company_id = ?", filter.CompanyID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	// Ordering
	orderBy := "created_at"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	order := "DESC"
	if filter.Order == "asc" {
		order = "ASC"
	}

	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Execute query
	if err := query.Limit(filter.Limit).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// GetByCompanyID retrieves jobs by company ID
func (r *JobRepository) GetByCompanyID(companyID uint) ([]*models.Job, error) {
	var jobs []*models.Job
	if err := r.db.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetByUserID retrieves jobs posted by a user
func (r *JobRepository) GetByUserID(userID uint) ([]*models.Job, error) {
	var jobs []*models.Job
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// IncrementViews increments job views count
func (r *JobRepository) IncrementViews(jobID uint) error {
	return r.db.Model(&models.Job{}).Where("id = ?", jobID).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// IncrementApplications increments job applications count
func (r *JobRepository) IncrementApplications(jobID uint) error {
	return r.db.Model(&models.Job{}).Where("id = ?", jobID).UpdateColumn("applications", gorm.Expr("applications + ?", 1)).Error
}

// DecrementApplications decrements job applications count
func (r *JobRepository) DecrementApplications(jobID uint) error {
	return r.db.Model(&models.Job{}).Where("id = ?", jobID).UpdateColumn("applications", gorm.Expr("applications - ?", 1)).Error
}

// GetPopularJobs retrieves popular jobs (most views)
func (r *JobRepository) GetPopularJobs(limit int) ([]*models.Job, error) {
	var jobs []*models.Job
	if err := r.db.Where("status = ?", models.JobStatusPublished).
		Order("views DESC").
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetRecentJobs retrieves recent jobs
func (r *JobRepository) GetRecentJobs(limit int) ([]*models.Job, error) {
	var jobs []*models.Job
	if err := r.db.Where("status = ?", models.JobStatusPublished).
		Order("created_at DESC").
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}
