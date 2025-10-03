package services

import (
	"errors"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
)

type CertificationService interface {
	Create(userID uint, req *models.CertificationRequest) (*models.Certification, error)
	GetAll(userID uint) ([]models.Certification, error)
	GetByID(userID uint, id uint) (*models.Certification, error)
	Update(userID uint, id uint, req *models.CertificationRequest) (*models.Certification, error)
	Delete(userID uint, id uint) error
}

type certificationService struct {
	repo           repository.CertificationRepository
	profileService ProfileService
}

func NewCertificationService(repo repository.CertificationRepository, profileService ProfileService) CertificationService {
	return &certificationService{
		repo:           repo,
		profileService: profileService,
	}
}

func (s *certificationService) Create(userID uint, req *models.CertificationRequest) (*models.Certification, error) {
	// Get or auto-create profile if not exists
	profile, err := s.profileService.GetOrCreateProfile(userID)
	if err != nil {
		return nil, errors.New("failed to get or create profile: " + err.Error())
	}

	certification := &models.Certification{
		ProfileID:         profile.ID,
		CertificationName: req.CertificationName,
		Organizer:         req.Organizer,
		IssueDate:         *req.IssueDate.ToTime(),
		ExpiryDate:        req.IssueDate.ToTime(),
		CredentialID:      req.CredentialID,
		CredentialURL:     req.CredentialURL,
		Description:       req.Description,
	}

	err = s.repo.Create(certification)
	if err != nil {
		return nil, err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return certification, nil
}

func (s *certificationService) GetAll(userID uint) ([]models.Certification, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	return s.repo.GetByProfileID(profile.ID)
}

func (s *certificationService) GetByID(userID uint, id uint) (*models.Certification, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	certification, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("certification not found")
	}

	if certification.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	return certification, nil
}

func (s *certificationService) Update(userID uint, id uint, req *models.CertificationRequest) (*models.Certification, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	certification, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("certification not found")
	}

	if certification.ProfileID != profile.ID {
		return nil, errors.New("unauthorized access")
	}

	certification.CertificationName = req.CertificationName
	certification.Organizer = req.Organizer
	certification.IssueDate = *req.IssueDate.ToTime()
	certification.ExpiryDate = req.ExpiryDate.ToTime()
	certification.CredentialID = req.CredentialID
	certification.CredentialURL = req.CredentialURL
	certification.Description = req.Description

	err = s.repo.Update(certification)
	if err != nil {
		return nil, err
	}

	return certification, nil
}

func (s *certificationService) Delete(userID uint, id uint) error {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return errors.New("profile not found")
	}

	certification, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("certification not found")
	}

	if certification.ProfileID != profile.ID {
		return errors.New("unauthorized access")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return nil
}
