package repository

import (
	"jobfair-user-profile-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(profile *models.Profile) error
	GetByID(id uuid.UUID) (*models.Profile, error)
	GetByUserID(userID uuid.UUID) (*models.Profile, error)
	Update(profile *models.Profile) error
	Delete(id uuid.UUID) error
	GetWithRelations(userID uuid.UUID) (*models.Profile, error)
	UpdateCompletionStatus(profileID uuid.UUID, status int) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) Create(profile *models.Profile) error {
	return r.db.Create(profile).Error
}

func (r *profileRepository) GetByID(id uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("id = ?", id).First(&profile).Error
	return &profile, err
}

func (r *profileRepository) GetByUserID(userID uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (r *profileRepository) Update(profile *models.Profile) error {
	return r.db.Save(profile).Error
}

func (r *profileRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Profile{}, "id = ?", id).Error
}

func (r *profileRepository) GetWithRelations(userID uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.
		Preload("WorkExperiences").
		Preload("Educations").
		Preload("Certifications").
		Preload("Skills").
		Preload("PositionPreferences").
		Preload("CareerPreference").
		Preload("CVDocument").
		Preload("Badges").
		Where("user_id = ?", userID).
		First(&profile).Error
	return &profile, err
}

func (r *profileRepository) UpdateCompletionStatus(profileID uuid.UUID, status int) error {
	return r.db.Model(&models.Profile{}).
		Where("id = ?", profileID).
		Update("completion_status", status).Error
}
