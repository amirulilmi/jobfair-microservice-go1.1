package services

import (
	"errors"
	"log"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
)

type PreferenceService interface {
	CreateOrUpdateCareerPreference(userID uint, req *models.CareerPreferenceRequest) (*models.CareerPreference, error)
	GetCareerPreference(userID uint) (*models.CareerPreference, error)
	CreatePositionPreferences(userID uint, positions []string) ([]models.PositionPreference, error)
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

func (s *preferenceService) CreateOrUpdateCareerPreference(userID uint, req *models.CareerPreferenceRequest) (*models.CareerPreference, error) {
	// Get or auto-create profile if not exists
	profile, err := s.profileService.GetOrCreateProfile(userID)
	if err != nil {
		return nil, errors.New("failed to get or create profile: " + err.Error())
	}

	log.Printf("[CreateCareerPreference] Got profile: ID=%d, UserID=%d", profile.ID, profile.UserID)

	// Map the request fields to the database model
	// job_type and work_location will be combined into PreferredWorkTypes
	workTypes := req.JobType
	if req.WorkLocation != "" {
		if workTypes != "" {
			workTypes += "," + req.WorkLocation
		} else {
			workTypes = req.WorkLocation
		}
	}

	// Check if preference already exists
	existing, _ := s.repo.GetCareerPreferenceByProfileID(profile.ID)

	if existing != nil {
		log.Printf("[CreateCareerPreference] Updating existing preference ID=%d", existing.ID)
		// Update existing
		existing.ExpectedSalaryMin = req.ExpectedSalaryMin
		existing.ExpectedSalaryMax = req.ExpectedSalaryMax
		existing.SalaryCurrency = req.Currency
		existing.PreferredWorkTypes = workTypes
		existing.WillingToRelocate = req.WillingToRelocate
		existing.AvailableStartDate = *req.AvailableFrom.ToTime()
		existing.IsActivelyLooking = true // Default to true when updating

		err = s.repo.UpdateCareerPreference(existing)
		if err != nil {
			return nil, err
		}

		s.profileService.UpdateCompletionStatus(userID)
		return existing, nil
	}

	log.Printf("[CreateCareerPreference] Creating new preference for ProfileID=%d", profile.ID)
	
	// Create new
	preference := &models.CareerPreference{
		ProfileID:          profile.ID,
		IsActivelyLooking:  true,
		ExpectedSalaryMin:  req.ExpectedSalaryMin,
		ExpectedSalaryMax:  req.ExpectedSalaryMax,
		SalaryCurrency:     req.Currency,
		IsNegotiable:       true, // Default to true
		PreferredWorkTypes: workTypes,
		WillingToRelocate:  req.WillingToRelocate,
		AvailableStartDate: *req.AvailableFrom.ToTime(),
	}

	log.Printf("[CreateCareerPreference] About to insert: ProfileID=%d, Salary=%d-%d", 
		preference.ProfileID, preference.ExpectedSalaryMin, preference.ExpectedSalaryMax)

	err = s.repo.CreateCareerPreference(preference)
	if err != nil {
		log.Printf("[CreateCareerPreference] Insert failed: %v", err)
		return nil, err
	}

	log.Printf("[CreateCareerPreference] Successfully created preference ID=%d", preference.ID)
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

func (s *preferenceService) CreatePositionPreferences(userID uint, positions []string) ([]models.PositionPreference, error) {
	// Get or auto-create profile if not exists
	profile, err := s.profileService.GetOrCreateProfile(userID)
	if err != nil {
		return nil, errors.New("failed to get or create profile: " + err.Error())
	}

	// Delete existing preferences
	_ = s.repo.DeletePositionPreferencesByProfileID(profile.ID)

	// Create new preferences from string array
	var preferences []models.PositionPreference
	for i, positionName := range positions {
		pref := models.PositionPreference{
			ProfileID:    profile.ID,
			PositionName: positionName,
			Priority:     i + 1, // Priority based on array order
		}

		err := s.repo.CreatePositionPreference(&pref)
		if err != nil {
			continue
		}

		preferences = append(preferences, pref)
	}

	if len(preferences) == 0 {
		return nil, errors.New("failed to create any position preferences")
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
