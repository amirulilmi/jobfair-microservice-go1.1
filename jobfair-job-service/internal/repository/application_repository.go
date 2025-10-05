package repository

import (
	"jobfair-job-service/internal/models"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// Create creates a new job application
func (r *ApplicationRepository) Create(application *models.JobApplication) (*models.JobApplication, error) {
	if err := r.db.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

// GetByID retrieves an application by ID
func (r *ApplicationRepository) GetByID(id uint) (*models.JobApplication, error) {
	var application models.JobApplication
	if err := r.db.Preload("Job").First(&application, id).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

// GetByJobIDAndUserID checks if user already applied to a job
func (r *ApplicationRepository) GetByJobIDAndUserID(jobID, userID uint) (*models.JobApplication, error) {
	var application models.JobApplication
	if err := r.db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

// Update updates an application
func (r *ApplicationRepository) Update(application *models.JobApplication) error {
	return r.db.Save(application).Error
}

// Delete soft deletes an application
func (r *ApplicationRepository) Delete(id uint) error {
	return r.db.Delete(&models.JobApplication{}, id).Error
}

// GetByUserID retrieves all applications by a user
func (r *ApplicationRepository) GetByUserID(userID uint, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.Model(&models.JobApplication{}).Where("user_id = ?", userID).Preload("Job")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// GetByJobID retrieves all applications for a job
func (r *ApplicationRepository) GetByJobID(jobID uint, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.Model(&models.JobApplication{}).Where("job_id = ?", jobID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// GetByCompanyID retrieves all applications for a company's jobs
func (r *ApplicationRepository) GetByCompanyID(companyID uint, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.Model(&models.JobApplication{}).
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.company_id = ?", companyID).
		Preload("Job")

	if status != "" {
		query = query.Where("job_applications.status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("job_applications.created_at DESC").Limit(limit).Offset(offset).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// CountByJobID counts applications for a job
func (r *ApplicationRepository) CountByJobID(jobID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.JobApplication{}).Where("job_id = ?", jobID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByUserID counts applications by a user
func (r *ApplicationRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.JobApplication{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatus counts applications by status
func (r *ApplicationRepository) CountByStatus(jobID uint, status models.ApplicationStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&models.JobApplication{}).
		Where("job_id = ? AND status = ?", jobID, status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetApplicationStats retrieves application statistics
func (r *ApplicationRepository) GetApplicationStats(companyID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	var results []struct {
		Status models.ApplicationStatus
		Count  int64
	}

	if err := r.db.Model(&models.JobApplication{}).
		Select("status, COUNT(*) as count").
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Where("jobs.company_id = ?", companyID).
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, result := range results {
		stats[string(result.Status)] = result.Count
	}

	return stats, nil
}

// GetByUserIDAndJobIDs checks which jobs are applied by user (batch query)
func (r *ApplicationRepository) GetByUserIDAndJobIDs(userID uint, jobIDs []uint) ([]*models.JobApplication, error) {
	var applications []*models.JobApplication
	if err := r.db.Where("user_id = ? AND job_id IN ?", userID, jobIDs).Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}
