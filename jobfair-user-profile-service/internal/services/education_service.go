package services

import (
	"errors"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
)

type EducationService interface {
	Create(userID uint, req *models.EducationRequest) (*models.Education, error)
	GetAll(userID uint) ([]models.Education, error)
	GetByID(userID uint, id uint) (*models.Education, error)
	Update(userID uint, id uint, req *models.EducationRequest) (*models.Education, error)
	Delete(userID uint, id uint) error
}

type educationService struct {
	repo           repository.EducationRepository
	profileService ProfileService
}

func NewEducationService(repo repository.EducationRepository, profileService ProfileService) EducationService {
	return &educationService{
		repo:           repo,
		profileService: profileService,
	}
}

func (s *educationService) Create(userID uint, req *models.EducationRequest) (*models.Education, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	education := &models.Education{
		ProfileID:   profile.ID,
		University:  req.University,
		Major:       req.Major,
		Degree:      req.Degree,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsCurrent:   req.IsCurrent,
		GPA:         req.GPA,
		Description: req.Description,
	}

	err = s.repo.Create(education)
	if err != nil {
		return nil, err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return education, nil
}

func (s *educationService) GetAll(userID uint) ([]models.Education, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	return s.repo.GetByProfileID(profile.ID)
}

func (s *educationService) GetByID(userID uint, id uint) (*models.Education, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	education, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("education not found")
	}

	if education.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	return education, nil
}

func (s *educationService) Update(userID uint, id uint, req *models.EducationRequest) (*models.Education, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	education, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("education not found")
	}

	if education.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	education.University = req.University
	education.Major = req.Major
	education.Degree = req.Degree
	education.StartDate = req.StartDate
	education.EndDate = req.EndDate
	education.IsCurrent = req.IsCurrent
	education.GPA = req.GPA
	education.Description = req.Description

	err = s.repo.Update(education)
	if err != nil {
		return nil, err
	}

	return education, nil
}

func (s *educationService) Delete(userID uint, id uint) error {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return errors.New("profile not found")
	}

	education, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("education not found")
	}

	if education.ProfileID != profile.ID {
		return errors.New("unauthorized access")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return nil
}
