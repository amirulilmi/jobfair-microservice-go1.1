package services

import (
	"errors"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
)

type SkillService interface {
	Create(userID uint, req *models.SkillRequest) (*models.Skill, error)
	CreateBulk(userID uint, req *models.BulkSkillRequest) ([]models.Skill, error)
	GetAll(userID uint) ([]models.Skill, error)
	GetByID(userID uint, id uint) (*models.Skill, error)
	Update(userID uint, id uint, req *models.SkillRequest) (*models.Skill, error)
	Delete(userID uint, id uint) error
}

type skillService struct {
	repo           repository.SkillRepository
	profileService ProfileService
}

func NewSkillService(repo repository.SkillRepository, profileService ProfileService) SkillService {
	return &skillService{
		repo:           repo,
		profileService: profileService,
	}
}

func (s *skillService) Create(userID uint, req *models.SkillRequest) (*models.Skill, error) {
	// Get or auto-create profile if not exists
	profile, err := s.profileService.GetOrCreateProfile(userID)
	if err != nil {
		return nil, errors.New("failed to get or create profile: " + err.Error())
	}

	skill := &models.Skill{
		ProfileID:         profile.ID,
		SkillName:         req.SkillName,
		SkillType:         req.SkillType,
		ProficiencyLevel:  req.ProficiencyLevel,
		YearsOfExperience: req.YearsOfExperience,
	}

	err = s.repo.Create(skill)
	if err != nil {
		return nil, err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return skill, nil
}

func (s *skillService) CreateBulk(userID uint, req *models.BulkSkillRequest) ([]models.Skill, error) {
	// Get or auto-create profile if not exists
	profile, err := s.profileService.GetOrCreateProfile(userID)
	if err != nil {
		return nil, errors.New("failed to get or create profile: " + err.Error())
	}

	var skills []models.Skill

	// Create all skills from the array
	for _, skillReq := range req.Skills {
		skill := models.Skill{
			ProfileID:         profile.ID,
			SkillName:         skillReq.SkillName,
			SkillType:         skillReq.SkillType,
			ProficiencyLevel:  skillReq.ProficiencyLevel,
			YearsOfExperience: skillReq.YearsOfExperience,
		}
		err := s.repo.Create(&skill)
		if err == nil {
			skills = append(skills, skill)
		}
	}

	if len(skills) == 0 {
		return nil, errors.New("failed to create any skills")
	}

	s.profileService.UpdateCompletionStatus(userID)
	return skills, nil
}

func (s *skillService) GetAll(userID uint) ([]models.Skill, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	return s.repo.GetByProfileID(profile.ID)
}

func (s *skillService) GetByID(userID uint, id uint) (*models.Skill, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	skill, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("skill not found")
	}

	if skill.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	return skill, nil
}

func (s *skillService) Update(userID uint, id uint, req *models.SkillRequest) (*models.Skill, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	skill, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("skill not found")
	}

	if skill.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	skill.SkillName = req.SkillName
	skill.SkillType = req.SkillType
	skill.ProficiencyLevel = req.ProficiencyLevel
	skill.YearsOfExperience = req.YearsOfExperience

	err = s.repo.Update(skill)
	if err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *skillService) Delete(userID uint, id uint) error {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return errors.New("profile not found")
	}

	skill, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("skill not found")
	}

	if skill.ProfileID != profile.ID {
		return errors.New("unauthorized access")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return nil
}
