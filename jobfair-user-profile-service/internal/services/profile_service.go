package services

import (
	"errors"
	"log"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"

	"gorm.io/gorm"
)

type ProfileService interface {
	CreateProfile(userID uint, fullName, phoneNumber string) (*models.Profile, error)
	GetProfile(userID uint) (*models.Profile, error)
	GetOrCreateProfile(userID uint) (*models.Profile, error)
	GetProfileWithRelations(userID uint) (*models.Profile, error)
	UpdateProfile(userID uint, req *models.ProfileUpdateRequest) (*models.Profile, error)
	CalculateCompletionStatus(profile *models.Profile) int
	UpdateCompletionStatus(userID uint) error
}

type profileService struct {
	profileRepo       repository.ProfileRepository
	workExpRepo       repository.WorkExperienceRepository
	educationRepo     repository.EducationRepository
	certificationRepo repository.CertificationRepository
	skillRepo         repository.SkillRepository
	preferenceRepo    repository.PreferenceRepository
	cvRepo            repository.CVRepository
}

func NewProfileService(
	profileRepo repository.ProfileRepository,
	workExpRepo repository.WorkExperienceRepository,
	educationRepo repository.EducationRepository,
	certificationRepo repository.CertificationRepository,
	skillRepo repository.SkillRepository,
	preferenceRepo repository.PreferenceRepository,
	cvRepo repository.CVRepository,
) ProfileService {
	return &profileService{
		profileRepo:       profileRepo,
		workExpRepo:       workExpRepo,
		educationRepo:     educationRepo,
		certificationRepo: certificationRepo,
		skillRepo:         skillRepo,
		preferenceRepo:    preferenceRepo,
		cvRepo:            cvRepo,
	}
}

func (s *profileService) CreateProfile(userID uint, fullName, phoneNumber string) (*models.Profile, error) {
	profile := &models.Profile{
		UserID:           userID,
		FullName:         fullName,
		PhoneNumber:      phoneNumber,
		CompletionStatus: 0,
	}

	err := s.profileRepo.Create(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *profileService) GetProfile(userID uint) (*models.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return profile, nil
}

func (s *profileService) GetOrCreateProfile(userID uint) (*models.Profile, error) {
	log.Printf("[GetOrCreateProfile] Starting for userID: %d", userID)
	
	// Try to get existing profile
	profile, err := s.profileRepo.GetByUserID(userID)
	if err == nil {
		log.Printf("[GetOrCreateProfile] Found existing profile: ID=%d, UserID=%d", profile.ID, profile.UserID)
		return profile, nil
	}

	log.Printf("[GetOrCreateProfile] Profile not found, error: %v", err)

	// If profile not found, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[GetOrCreateProfile] Creating new profile for userID: %d", userID)
		
		newProfile := &models.Profile{
			UserID:           userID,
			FullName:         "",
			PhoneNumber:      "",
			CompletionStatus: 0,
		}

		err = s.profileRepo.Create(newProfile)
		if err != nil {
			log.Printf("[GetOrCreateProfile] Failed to create profile: %v", err)
			return nil, err
		}
		
		// GORM should auto-populate the ID after Create
		log.Printf("[GetOrCreateProfile] Profile created with ID=%d", newProfile.ID)
		
		// Verify ID is set
		if newProfile.ID == 0 {
			log.Printf("[GetOrCreateProfile] WARNING: ID is 0 after create, reloading from DB")
			// Reload from DB to ensure ID is populated
			profile, err = s.profileRepo.GetByUserID(userID)
			if err != nil {
				log.Printf("[GetOrCreateProfile] Failed to reload profile: %v", err)
				return nil, err
			}
			log.Printf("[GetOrCreateProfile] Profile reloaded: ID=%d", profile.ID)
			return profile, nil
		}
		
		return newProfile, nil
	}

	log.Printf("[GetOrCreateProfile] Unexpected error: %v", err)
	return nil, err
}

func (s *profileService) GetProfileWithRelations(userID uint) (*models.Profile, error) {
	profile, err := s.profileRepo.GetWithRelations(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return profile, nil
}

func (s *profileService) UpdateProfile(userID uint, req *models.ProfileUpdateRequest) (*models.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	// Update only provided fields
	if req.FirstName != nil {
		profile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profile.LastName = *req.LastName
	}
	if req.FullName != nil {
		profile.FullName = *req.FullName
	}
	
	// Auto-generate full_name if first_name or last_name is provided but full_name is not
	if (req.FirstName != nil || req.LastName != nil) && req.FullName == nil {
		// Build full name from first_name and last_name
		firstName := profile.FirstName
		lastName := profile.LastName
		
		if firstName != "" && lastName != "" {
			profile.FullName = firstName + " " + lastName
		} else if firstName != "" {
			profile.FullName = firstName
		} else if lastName != "" {
			profile.FullName = lastName
		}
	}
	
	if req.PhoneNumber != nil {
		profile.PhoneNumber = *req.PhoneNumber
	}
	if req.Headline != nil {
		profile.Headline = *req.Headline
	}
	if req.Summary != nil {
		profile.Summary = *req.Summary
	}
	if req.Bio != nil {
		profile.Bio = *req.Bio
	}
	if req.Location != nil {
		profile.Location = *req.Location
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != nil {
		profile.Gender = *req.Gender
	}
	if req.Address != nil {
		profile.Address = *req.Address
	}
	if req.City != nil {
		profile.City = *req.City
	}
	if req.Province != nil {
		profile.Province = *req.Province
	}
	if req.Country != nil {
		profile.Country = *req.Country
	}
	if req.PostalCode != nil {
		profile.PostalCode = *req.PostalCode
	}
	if req.LinkedInURL != nil {
		profile.LinkedInURL = *req.LinkedInURL
	}
	if req.GitHubURL != nil {
		profile.GitHubURL = *req.GitHubURL
	}
	if req.PortfolioURL != nil {
		profile.PortfolioURL = *req.PortfolioURL
	}
	if req.ProfilePictureURL != nil {
		profile.ProfilePictureURL = *req.ProfilePictureURL
	}
	if req.BannerImageURL != nil {
		profile.BannerImageURL = *req.BannerImageURL
	}

	err = s.profileRepo.Update(profile)
	if err != nil {
		return nil, err
	}

	// Update completion status
	s.UpdateCompletionStatus(userID)

	return profile, nil
}

func (s *profileService) CalculateCompletionStatus(profile *models.Profile) int {
	totalFields := 15
	completedFields := 0

	// Basic info (5 fields)
	if profile.FullName != "" {
		completedFields++
	}
	if profile.PhoneNumber != "" {
		completedFields++
	}
	if profile.Bio != "" {
		completedFields++
	}
	if profile.DateOfBirth != nil {
		completedFields++
	}
	if profile.City != "" {
		completedFields++
	}

	// Work experience (1 field)
	workExps, _ := s.workExpRepo.GetByProfileID(profile.ID)
	if len(workExps) > 0 {
		completedFields++
	}

	// Education (1 field)
	educations, _ := s.educationRepo.GetByProfileID(profile.ID)
	if len(educations) > 0 {
		completedFields++
	}

	// Certifications (1 field)
	certifications, _ := s.certificationRepo.GetByProfileID(profile.ID)
	if len(certifications) > 0 {
		completedFields++
	}

	// Skills (2 fields - technical and soft)
	technicalSkills, _ := s.skillRepo.GetByProfileIDAndType(profile.ID, "technical")
	if len(technicalSkills) > 0 {
		completedFields++
	}
	softSkills, _ := s.skillRepo.GetByProfileIDAndType(profile.ID, "soft")
	if len(softSkills) > 0 {
		completedFields++
	}

	// Career preference (1 field)
	careerPref, _ := s.preferenceRepo.GetCareerPreferenceByProfileID(profile.ID)
	if careerPref != nil && careerPref.ID != 0 {
		completedFields++
	}

	// Position preferences (1 field)
	positionPrefs, _ := s.preferenceRepo.GetPositionPreferencesByProfileID(profile.ID)
	if len(positionPrefs) > 0 {
		completedFields++
	}

	// CV Document (1 field)
	cv, _ := s.cvRepo.GetByProfileID(profile.ID)
	if cv != nil && cv.ID != 0 {
		completedFields++
	}

	// Profile picture (1 field)
	if profile.ProfilePictureURL != "" {
		completedFields++
	}

	// LinkedIn (1 field)
	if profile.LinkedInURL != "" {
		completedFields++
	}

	// Calculate percentage
	completionPercentage := (completedFields * 100) / totalFields
	return completionPercentage
}

func (s *profileService) UpdateCompletionStatus(userID uint) error {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	completionStatus := s.CalculateCompletionStatus(profile)
	return s.profileRepo.UpdateCompletionStatus(profile.ID, completionStatus)
}
