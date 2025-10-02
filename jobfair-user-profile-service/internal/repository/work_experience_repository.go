package repository

import (
	"jobfair-user-profile-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkExperienceRepository interface {
	Create(workExp *models.WorkExperience) error
	GetByID(id uuid.UUID) (*models.WorkExperience, error)
	GetByProfileID(profileID uuid.UUID) ([]models.WorkExperience, error)
	Update(workExp *models.WorkExperience) error
	Delete(id uuid.UUID) error
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

func (r *workExperienceRepository) GetByID(id uuid.UUID) (*models.WorkExperience, error) {
	var workExp models.WorkExperience
	err := r.db.Where("id = ?", id).First(&workExp).Error
	return &workExp, err
}

func (r *workExperienceRepository) GetByProfileID(profileID uuid.UUID) ([]models.WorkExperience, error) {
	var workExps []models.WorkExperience
	err := r.db.Where("profile_id = ?", profileID).Order("start_date DESC").Find(&workExps).Error
	return workExps, err
}

func (r *workExperienceRepository) Update(workExp *models.WorkExperience) error {
	return r.db.Save(workExp).Error
}

func (r *workExperienceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.WorkExperience{}, "id = ?", id).Error
}
