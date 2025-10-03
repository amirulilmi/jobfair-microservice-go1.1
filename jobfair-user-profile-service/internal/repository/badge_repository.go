package repository

import (
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type BadgeRepository interface {
	Create(badge *models.Badge) error
	GetByID(id uint) (*models.Badge, error)
	GetAll() ([]models.Badge, error)
	Update(badge *models.Badge) error
	Delete(id uint) error
	
	// Profile Badge operations
	AssignBadgeToProfile(profileID, badgeID uint) error
	GetProfileBadges(profileID uint) ([]models.Badge, error)
	RemoveBadgeFromProfile(profileID, badgeID uint) error
}

type badgeRepository struct {
	db *gorm.DB
}

func NewBadgeRepository(db *gorm.DB) BadgeRepository {
	return &badgeRepository{db: db}
}

func (r *badgeRepository) Create(badge *models.Badge) error {
	return r.db.Create(badge).Error
}

func (r *badgeRepository) GetByID(id uint) (*models.Badge, error) {
	var badge models.Badge
	err := r.db.Where("id = ?", id).First(&badge).Error
	return &badge, err
}

func (r *badgeRepository) GetAll() ([]models.Badge, error) {
	var badges []models.Badge
	err := r.db.Find(&badges).Error
	return badges, err
}

func (r *badgeRepository) Update(badge *models.Badge) error {
	return r.db.Save(badge).Error
}

func (r *badgeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Badge{}, "id = ?", id).Error
}

func (r *badgeRepository) AssignBadgeToProfile(profileID, badgeID uint) error {
	profileBadge := models.ProfileBadge{
		ProfileID: profileID,
		BadgeID:   badgeID,
	}
	return r.db.Create(&profileBadge).Error
}

func (r *badgeRepository) GetProfileBadges(profileID uint) ([]models.Badge, error) {
	var badges []models.Badge
	err := r.db.
		Joins("JOIN profile_badges ON profile_badges.badge_id = badges.id").
		Where("profile_badges.profile_id = ?", profileID).
		Find(&badges).Error
	return badges, err
}

func (r *badgeRepository) RemoveBadgeFromProfile(profileID, badgeID uint) error {
	return r.db.
		Where("profile_id = ? AND badge_id = ?", profileID, badgeID).
		Delete(&models.ProfileBadge{}).Error
}
