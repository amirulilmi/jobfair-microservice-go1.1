package repository

import (
	"jobfair-user-profile-service/internal/models"

	"gorm.io/gorm"
)

type SkillRepository interface {
	Create(skill *models.Skill) error
	GetByID(id uint) (*models.Skill, error)
	GetByProfileID(profileID uint) ([]models.Skill, error)
	GetByProfileIDAndType(profileID uint, skillType string) ([]models.Skill, error)
	Update(skill *models.Skill) error
	Delete(id uint) error
}

type skillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepository{db: db}
}

func (r *skillRepository) Create(skill *models.Skill) error {
	return r.db.Create(skill).Error
}

func (r *skillRepository) GetByID(id uint) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.Where("id = ?", id).First(&skill).Error
	return &skill, err
}

func (r *skillRepository) GetByProfileID(profileID uint) ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Where("profile_id = ?", profileID).Find(&skills).Error
	return skills, err
}

func (r *skillRepository) GetByProfileIDAndType(profileID uint, skillType string) ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Where("profile_id = ? AND skill_type = ?", profileID, skillType).Find(&skills).Error
	return skills, err
}

func (r *skillRepository) Update(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *skillRepository) Delete(id uint) error {
	return r.db.Delete(&models.Skill{}, "id = ?", id).Error
}
