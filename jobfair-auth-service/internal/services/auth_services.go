package services

import (
	"errors"
	// "time"

	"jobfair-auth-service/internal/models"
	"jobfair-auth-service/internal/repository"
	"jobfair-auth-service/internal/utils"
)

type AuthService struct {
	userRepo           *repository.UserRepository
	jobSeekerRepo      *repository.JobSeekerProfileRepository
	companyProfileRepo *repository.CompanyBasicProfileRepository
	jwtSecret          string
}

func NewAuthService(
	userRepo *repository.UserRepository,
	jobSeekerRepo *repository.JobSeekerProfileRepository,
	companyProfileRepo *repository.CompanyBasicProfileRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		jobSeekerRepo:      jobSeekerRepo,
		companyProfileRepo: companyProfileRepo,
		jwtSecret:          jwtSecret,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		UserType: req.UserType,
		IsActive: true,
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, string(user.UserType), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	userID, err := utils.ValidateRefreshToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, string(user.UserType), s.jwtSecret)
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// GetCurrentUserWithProfile returns complete user data with profile based on user type
func (s *AuthService) GetCurrentUserWithProfile(userID uint) (map[string]interface{}, error) {
	// Get user data
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Base response
	response := map[string]interface{}{
		"id":                  user.ID,
		"email":               user.Email,
		"user_type":           user.UserType,
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"phone_number":        user.PhoneNumber,
		"country_code":        user.CountryCode,
		"country":             user.Country,
		"profile_photo":       user.ProfilePhoto,
		"is_email_verified":   user.IsEmailVerified,
		"is_phone_verified":   user.IsPhoneVerified,
		"email_verified_at":   user.EmailVerifiedAt,
		"phone_verified_at":   user.PhoneVerifiedAt,
		"is_active":           user.IsActive,
		"is_profile_complete": user.IsProfileComplete,
		"created_at":          user.CreatedAt,
		"updated_at":          user.UpdatedAt,
	}

	// Add profile data based on user type
	switch user.UserType {
	case models.UserTypeJobSeeker:
		profile, err := s.jobSeekerRepo.GetByUserID(userID)
		if err == nil {
			response["profile"] = map[string]interface{}{
				"current_job_title":      profile.CurrentJobTitle,
				"current_company":        profile.CurrentCompany,
				"employment_status":      profile.EmploymentStatus,
				"job_search_status":      profile.JobSearchStatus,
				"desired_positions":      profile.DesiredPositions,
				"preferred_locations":    profile.PreferredLocations,
				"job_types":              profile.JobTypes,
				"notifications_enabled":  profile.NotificationsEnabled,
				"location_enabled":       profile.LocationEnabled,
			}
		}

	case models.UserTypeCompany:
		profile, err := s.companyProfileRepo.GetByUserID(userID)
		if err == nil {
			response["company_profile"] = map[string]interface{}{
				"company_name": profile.CompanyName,
				"industry":     profile.Industry,
				"contact_name": profile.ContactName,
				"phone_number": profile.PhoneNumber,
				"address":      profile.Address,
				"website":      profile.Website,
				"logo_url":     profile.LogoURL,
			}
		}
	}

	return response, nil
}
