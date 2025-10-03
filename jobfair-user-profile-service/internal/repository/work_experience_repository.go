package repository

import (
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type WorkExperienceRepository interface {
	Create(workExp *models.WorkExperience) error
	GetByID(id uint) (*models.WorkExperience, error)
	GetByProfileID(profileID uint) ([]models.WorkExperience, error)
	Update(workExp *models.WorkExperience) error
	Delete(id uint) error
}

type workExperienceRepository struct {
	db *gorm.DB
}

func NewWorkExperienceRepository(db *gorm.DB) WorkExperienceRepository {
	return &workExperienceRepository{db: db}
}

func (r *workExperienceRepository) Create(workExp *models.WorkExperience) error {
	return r.db.Create(workExp).Error
}

func (r *workExperienceRepository) GetByID(id uint) (*models.WorkExperience, error) {
	var workExp models.WorkExperience
	err := r.db.Where("id = ?", id).First(&workExp).Error
	return &workExp, err
}

func (r *workExperienceRepository) GetByProfileID(profileID uint) ([]models.WorkExperience, error) {
	var workExps []models.WorkExperience
	err := r.db.Where("profile_id = ?", profileID).Order("start_date DESC").Find(&workExps).Error
	return workExps, err
}

func (r *workExperienceRepository) Update(workExp *models.WorkExperience) error {
	return r.db.Save(workExp).Error
}

func (r *workExperienceRepository) Delete(id uint) error {
	return r.db.Delete(&models.WorkExperience{}, "id = ?", id).Error
}
