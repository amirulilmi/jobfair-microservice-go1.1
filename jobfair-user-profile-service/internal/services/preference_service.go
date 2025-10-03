package services

import (
	"errors"
	"time"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
)

type PreferenceService interface {
	CreateOrUpdateCareerPreference(userID uint, isActivelyLooking bool, salaryMin, salaryMax *int, currency string, isNegotiable bool, workTypes, locations string, relocate bool, startDate *time.Time) (*models.CareerPreference, error)
	GetCareerPreference(userID uint) (*models.CareerPreference, error)
	CreatePositionPreferences(userID uint, positions []models.PositionPreferenceRequest) ([]models.PositionPreference, error)
	GetPositionPreferences(userID uint) ([]models.PositionPreference, error)
	DeletePositionPreference(userID uint, id uint) error
}

type preferenceService struct {
	repo           repository.PreferenceRepository
	profileService ProfileService
}

func NewPreferenceService(repo repository.PreferenceRepository, profileService ProfileService) PreferenceService {
	return &preferenceService{
		repo:           repo,
		profileService: profileService,
	}
}

func (s *preferenceService) CreateOrUpdateCareerPreference(userID uint, isActivelyLooking bool, salaryMin, salaryMax *int, currency string, isNegotiable bool, workTypes, locations string, relocate bool, startDate *time.Time) (*models.CareerPreference, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	// Check if preference already exists
	existing, _ := s.repo.GetCareerPreferenceByProfileID(profile.ID)
	
	if existing != nil {
		// Update existing
		existing.IsActivelyLooking = isActivelyLooking
		existing.ExpectedSalaryMin = salaryMin
		existing.ExpectedSalaryMax = salaryMax
		existing.SalaryCurrency = currency
		existing.IsNegotiable = isNegotiable
		existing.PreferredWorkTypes = workTypes
		existing.PreferredLocations = locations
		existing.WillingToRelocate = relocate
		existing.AvailableStartDate = startDate
		
		err = s.repo.UpdateCareerPreference(existing)
		if err != nil {
			return nil, err
		}
		
		s.profileService.UpdateCompletionStatus(userID)
		return existing, nil
	}

	// Create new
	preference := &models.CareerPreference{
		ProfileID:          profile.ID,
		IsActivelyLooking:  isActivelyLooking,
		ExpectedSalaryMin:  salaryMin,
		ExpectedSalaryMax:  salaryMax,
		SalaryCurrency:     currency,
		IsNegotiable:       isNegotiable,
		PreferredWorkTypes: workTypes,
		PreferredLocations: locations,
		WillingToRelocate:  relocate,
		AvailableStartDate: startDate,
	}

	err = s.repo.CreateCareerPreference(preference)
	if err != nil {
		return nil, err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return preference, nil
}

func (s *preferenceService) GetCareerPreference(userID uint) (*models.CareerPreference, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	preference, err := s.repo.GetCareerPreferenceByProfileID(profile.ID)
	if err != nil {
		return nil, errors.New("career preference not found")
	}

	return preference, nil
}

func (s *preferenceService) CreatePositionPreferences(userID uint, positions []models.PositionPreferenceRequest) ([]models.PositionPreference, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	// Delete existing preferences
	_ = s.repo.DeletePositionPreferencesByProfileID(profile.ID)

	// Create new preferences
	var preferences []models.PositionPreference
	for _, pos := range positions {
		pref := models.PositionPreference{
			ProfileID:    profile.ID,
			PositionName: pos.PositionName,
			Priority:     pos.Priority,
		}
		
		err := s.repo.CreatePositionPreference(&pref)
		if err != nil {
			continue
		}
		
		preferences = append(preferences, pref)
	}

	s.profileService.UpdateCompletionStatus(userID)
	return preferences, nil
}

func (s *preferenceService) GetPositionPreferences(userID uint) ([]models.PositionPreference, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	return s.repo.GetPositionPreferencesByProfileID(profile.ID)
}

func (s *preferenceService) DeletePositionPreference(userID uint, id uint) error {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return errors.New("profile not found")
	}

	pref, err := s.repo.GetPositionPreferenceByID(id)
	if err != nil {
		return errors.New("position preference not found")
	}

	if pref.ProfileID != profile.ID {
		return errors.New("unauthorized access")
	}

	err = s.repo.DeletePositionPreference(id)
	if err != nil {
		return err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return nil
}
