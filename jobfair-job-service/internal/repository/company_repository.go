package repository

import (
	"time"

	"gorm.io/gorm"
)

type CompanyMapping struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;uniqueIndex"`
	CompanyID   uint      `gorm:"not null"`
	CompanyName string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for CompanyMapping
func (CompanyMapping) TableName() string {
	return "company_mappings"
}

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// GetCompanyIDByUserID gets company ID from user ID via mapping table
func (r *CompanyRepository) GetCompanyIDByUserID(userID uint) (uint, error) {
	var mapping CompanyMapping
	
	err := r.db.Table("company_mappings").
		Where("user_id = ?", userID).
		First(&mapping).Error
		
	if err != nil {
		return 0, err
	}
	
	return mapping.CompanyID, nil
}

// UpsertCompanyMapping creates or updates company mapping
func (r *CompanyRepository) UpsertCompanyMapping(userID, companyID uint, companyName string) error {
	mapping := CompanyMapping{
		UserID:      userID,
		CompanyID:   companyID,
		CompanyName: companyName,
		UpdatedAt:   time.Now(),
	}
	
	return r.db.Where("user_id = ?", userID).
		Assign(mapping).
		FirstOrCreate(&mapping).Error
}

// DeleteMappingByUserID deletes company mapping by user ID
func (r *CompanyRepository) DeleteMappingByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&CompanyMapping{}).Error
}
