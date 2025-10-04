package repository

import (
	"log"
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(profile *models.Profile) error
	GetByID(id uint) (*models.Profile, error)
	GetByUserID(userID uint) (*models.Profile, error)
	Update(profile *models.Profile) error
	Delete(id uint) error
	GetWithRelations(userID uint) (*models.Profile, error)
	UpdateCompletionStatus(profileID uint, status int) error
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

func (r *profileRepository) GetByID(id uint) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("id = ?", id).First(&profile).Error
	return &profile, err
}

func (r *profileRepository) GetByUserID(userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	// Debug log
	log.Printf("[ProfileRepo] GetByUserID(%d) found profile: ID=%d", userID, profile.ID)
	return &profile, nil
}

func (r *profileRepository) Update(profile *models.Profile) error {
	return r.db.Save(profile).Error
}

func (r *profileRepository) Delete(id uint) error {
	return r.db.Delete(&models.Profile{}, "id = ?", id).Error
}

func (r *profileRepository) GetWithRelations(userID uint) (*models.Profile, error) {
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

func (r *profileRepository) UpdateCompletionStatus(profileID uint, status int) error {
	return r.db.Model(&models.Profile{}).
		Where("id = ?", profileID).
		Update("completion_status", status).Error
}
