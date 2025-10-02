package repository

import (
	"jobfair-user-profile-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillRepository interface {
	Create(skill *models.Skill) error
	GetByID(id uuid.UUID) (*models.Skill, error)
	GetByProfileID(profileID uuid.UUID) ([]models.Skill, error)
	GetByProfileIDAndType(profileID uuid.UUID, skillType string) ([]models.Skill, error)
	Update(skill *models.Skill) error
	Delete(id uuid.UUID) error
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

func (r *skillRepository) GetByID(id uuid.UUID) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.Where("id = ?", id).First(&skill).Error
	return &skill, err
}

func (r *skillRepository) GetByProfileID(profileID uuid.UUID) ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Where("profile_id = ?", profileID).Find(&skills).Error
	return skills, err
}

func (r *skillRepository) GetByProfileIDAndType(profileID uuid.UUID, skillType string) ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Where("profile_id = ? AND skill_type = ?", profileID, skillType).Find(&skills).Error
	return skills, err
}

func (r *skillRepository) Update(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *skillRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Skill{}, "id = ?", id).Error
}
