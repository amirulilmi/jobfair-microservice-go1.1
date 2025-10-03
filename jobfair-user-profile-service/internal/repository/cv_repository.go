package repository

import (
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type CVRepository interface {
	Create(cv *models.CVDocument) error
	GetByID(id uint) (*models.CVDocument, error)
	GetByProfileID(profileID uint) (*models.CVDocument, error)
	Update(cv *models.CVDocument) error
	Delete(id uint) error
}

type cvRepository struct {
	db *gorm.DB
}

func NewCVRepository(db *gorm.DB) CVRepository {
	return &cvRepository{db: db}
}

func (r *cvRepository) Create(cv *models.CVDocument) error {
	return r.db.Create(cv).Error
}

func (r *cvRepository) GetByID(id uint) (*models.CVDocument, error) {
	var cv models.CVDocument
	err := r.db.Where("id = ?", id).First(&cv).Error
	return &cv, err
}

func (r *cvRepository) GetByProfileID(profileID uint) (*models.CVDocument, error) {
	var cv models.CVDocument
	err := r.db.Where("profile_id = ?", profileID).First(&cv).Error
	return &cv, err
}

func (r *cvRepository) Update(cv *models.CVDocument) error {
	return r.db.Save(cv).Error
}

func (r *cvRepository) Delete(id uint) error {
	return r.db.Delete(&models.CVDocument{}, "id = ?", id).Error
}
