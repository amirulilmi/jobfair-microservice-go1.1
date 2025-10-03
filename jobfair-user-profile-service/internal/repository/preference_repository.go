package repository

import (
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type PreferenceRepository interface {
	// Career Preference
	CreateCareerPreference(pref *models.CareerPreference) error
	GetCareerPreferenceByProfileID(profileID uint) (*models.CareerPreference, error)
	UpdateCareerPreference(pref *models.CareerPreference) error
	
	// Position Preference
	CreatePositionPreference(pref *models.PositionPreference) error
	GetPositionPreferenceByID(id uint) (*models.PositionPreference, error)
	GetPositionPreferencesByProfileID(profileID uint) ([]models.PositionPreference, error)
	DeletePositionPreference(id uint) error
	DeletePositionPreferencesByProfileID(profileID uint) error
}

type preferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) PreferenceRepository {
	return &preferenceRepository{db: db}
}

// Career Preference methods
func (r *preferenceRepository) CreateCareerPreference(pref *models.CareerPreference) error {
	return r.db.Create(pref).Error
}

func (r *preferenceRepository) GetCareerPreferenceByProfileID(profileID uint) (*models.CareerPreference, error) {
	var pref models.CareerPreference
	err := r.db.Where("profile_id = ?", profileID).First(&pref).Error
	return &pref, err
}

func (r *preferenceRepository) UpdateCareerPreference(pref *models.CareerPreference) error {
	return r.db.Save(pref).Error
}

// Position Preference methods
func (r *preferenceRepository) CreatePositionPreference(pref *models.PositionPreference) error {
	return r.db.Create(pref).Error
}

func (r *preferenceRepository) GetPositionPreferenceByID(id uint) (*models.PositionPreference, error) {
	var pref models.PositionPreference
	err := r.db.Where("id = ?", id).First(&pref).Error
	return &pref, err
}

func (r *preferenceRepository) GetPositionPreferencesByProfileID(profileID uint) ([]models.PositionPreference, error) {
	var prefs []models.PositionPreference
	err := r.db.Where("profile_id = ?", profileID).Order("priority ASC").Find(&prefs).Error
	return prefs, err
}

func (r *preferenceRepository) DeletePositionPreference(id uint) error {
	return r.db.Delete(&models.PositionPreference{}, "id = ?", id).Error
}

func (r *preferenceRepository) DeletePositionPreferencesByProfileID(profileID uint) error {
	return r.db.Delete(&models.PositionPreference{}, "profile_id = ?", profileID).Error
}
