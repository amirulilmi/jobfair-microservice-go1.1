package repository

import (
	"jobfair-job-service/internal/models"

	"gorm.io/gorm"
)

type SavedJobRepository struct {
	db *gorm.DB
}

func NewSavedJobRepository(db *gorm.DB) *SavedJobRepository {
	return &SavedJobRepository{db: db}
}

// Create saves a job
func (r *SavedJobRepository) Create(savedJob *models.SavedJob) (*models.SavedJob, error) {
	if err := r.db.Create(savedJob).Error; err != nil {
		return nil, err
	}
	return savedJob, nil
}

// GetByJobIDAndUserID checks if a job is saved by user
func (r *SavedJobRepository) GetByJobIDAndUserID(jobID, userID uint) (*models.SavedJob, error) {
	var savedJob models.SavedJob
	if err := r.db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&savedJob).Error; err != nil {
		return nil, err
	}
	return &savedJob, nil
}

// Delete unsaves a job
func (r *SavedJobRepository) Delete(id uint) error {
	return r.db.Delete(&models.SavedJob{}, id).Error
}

// DeleteByJobIDAndUserID unsaves a job by job ID and user ID
func (r *SavedJobRepository) DeleteByJobIDAndUserID(jobID, userID uint) error {
	return r.db.Where("job_id = ? AND user_id = ?", jobID, userID).Delete(&models.SavedJob{}).Error
}

// GetByUserID retrieves all saved jobs by a user
func (r *SavedJobRepository) GetByUserID(userID uint, page, limit int) ([]*models.SavedJob, int64, error) {
	var savedJobs []*models.SavedJob
	var total int64

	query := r.db.Model(&models.SavedJob{}).Where("user_id = ?", userID).Preload("Job")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&savedJobs).Error; err != nil {
		return nil, 0, err
	}

	return savedJobs, total, nil
}

// IsSaved checks if a job is saved by user
func (r *SavedJobRepository) IsSaved(jobID, userID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.SavedJob{}).
		Where("job_id = ? AND user_id = ?", jobID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByUserID counts saved jobs by a user
func (r *SavedJobRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.SavedJob{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
