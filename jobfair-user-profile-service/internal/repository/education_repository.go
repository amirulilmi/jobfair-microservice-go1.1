package repository

import (
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type EducationRepository interface {
	Create(education *models.Education) error
	GetByID(id uint) (*models.Education, error)
	GetByProfileID(profileID uint) ([]models.Education, error)
	Update(education *models.Education) error
	Delete(id uint) error
}

type educationRepository struct {
	db *gorm.DB
}

func NewEducationRepository(db *gorm.DB) EducationRepository {
	return &educationRepository{db: db}
}

func (r *educationRepository) Create(education *models.Education) error {
	return r.db.Create(education).Error
}

func (r *educationRepository) GetByID(id uint) (*models.Education, error) {
	var education models.Education
	err := r.db.Where("id = ?", id).First(&education).Error
	return &education, err
}

func (r *educationRepository) GetByProfileID(profileID uint) ([]models.Education, error) {
	var educations []models.Education
	err := r.db.Where("profile_id = ?", profileID).Order("start_date DESC").Find(&educations).Error
	return educations, err
}

func (r *educationRepository) Update(education *models.Education) error {
	return r.db.Save(education).Error
}

func (r *educationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Education{}, "id = ?", id).Error
}
