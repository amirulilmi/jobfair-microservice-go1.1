package services

import (
	"errors"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
)

type WorkExperienceService interface {
	Create(userID uint, req *models.WorkExperienceRequest) (*models.WorkExperience, error)
	GetAll(userID uint) ([]models.WorkExperience, error)
	GetByID(userID uint, id uint) (*models.WorkExperience, error)
	Update(userID uint, id uint, req *models.WorkExperienceRequest) (*models.WorkExperience, error)
	Delete(userID uint, id uint) error
}

type workExperienceService struct {
	workExpRepo    repository.WorkExperienceRepository
	profileService ProfileService
}

func NewWorkExperienceService(
	workExpRepo repository.WorkExperienceRepository,
	profileService ProfileService,
) WorkExperienceService {
	return &workExperienceService{
		workExpRepo:    workExpRepo,
		profileService: profileService,
	}
}

func (s *workExperienceService) Create(userID uint, req *models.WorkExperienceRequest) (*models.WorkExperience, error) {
	// Get profile first to ensure it exists and get profile ID
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	workExp := &models.WorkExperience{
		ProfileID:      profile.ID,
		CompanyName:    req.CompanyName,
		JobPosition:    req.JobPosition,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsCurrentJob:   req.IsCurrentJob,
		JobDescription: req.JobDescription,
	}

	err = s.workExpRepo.Create(workExp)
	if err != nil {
		return nil, err
	}

	// Update profile completion status
	s.profileService.UpdateCompletionStatus(userID)

	return workExp, nil
}

func (s *workExperienceService) GetAll(userID uint) ([]models.WorkExperience, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	return s.workExpRepo.GetByProfileID(profile.ID)
}

func (s *workExperienceService) GetByID(userID uint, id uint) (*models.WorkExperience, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	workExp, err := s.workExpRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("work experience not found")
	}

	// Verify ownership
	if workExp.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	return workExp, nil
}

func (s *workExperienceService) Update(userID uint, id uint, req *models.WorkExperienceRequest) (*models.WorkExperience, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	workExp, err := s.workExpRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("work experience not found")
	}

	// Verify ownership
	if workExp.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	// Update fields
	workExp.CompanyName = req.CompanyName
	workExp.JobPosition = req.JobPosition
	workExp.StartDate = req.StartDate
	workExp.EndDate = req.EndDate
	workExp.IsCurrentJob = req.IsCurrentJob
	workExp.JobDescription = req.JobDescription

	err = s.workExpRepo.Update(workExp)
	if err != nil {
		return nil, err
	}

	return workExp, nil
}

func (s *workExperienceService) Delete(userID uint, id uint) error {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return errors.New("profile not found")
	}

	workExp, err := s.workExpRepo.GetByID(id)
	if err != nil {
		return errors.New("work experience not found")
	}

	// Verify ownership
	if workExp.ProfileID != profile.ID {
		return errors.New("unauthorized access")
	}

	err = s.workExpRepo.Delete(id)
	if err != nil {
		return err
	}

	// Update profile completion status
	s.profileService.UpdateCompletionStatus(userID)

	return nil
}
