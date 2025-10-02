package repository

import (
	"jobfair-user-profile-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificationRepository interface {
	Create(certification *models.Certification) error
	GetByID(id uuid.UUID) (*models.Certification, error)
	GetByProfileID(profileID uuid.UUID) ([]models.Certification, error)
	Update(certification *models.Certification) error
	Delete(id uuid.UUID) error
}

type certificationRepository struct {
	db *gorm.DB
}

func NewCertificationRepository(db *gorm.DB) CertificationRepository {
	return &certificationRepository{db: db}
}

func (r *certificationRepository) Create(certification *models.Certification) error {
	return r.db.Create(certification).Error
}

func (r *certificationRepository) GetByID(id uuid.UUID) (*models.Certification, error) {
	var certification models.Certification
	err := r.db.Where("id = ?", id).First(&certification).Error
	return &certification, err
}

func (r *certificationRepository) GetByProfileID(profileID uuid.UUID) ([]models.Certification, error) {
	var certifications []models.Certification
	err := r.db.Where("profile_id = ?", profileID).Order("issue_date DESC").Find(&certifications).Error
	return certifications, err
}

func (r *certificationRepository) Update(certification *models.Certification) error {
	return r.db.Save(certification).Error
}

func (r *certificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Certification{}, "id = ?", id).Error
}
