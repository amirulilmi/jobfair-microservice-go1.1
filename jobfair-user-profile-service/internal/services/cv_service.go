package services

import (
	"errors"
	"fmt"
	"io"
	"jobfair-user-profile-service/internal/config"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/repository"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CVService interface {
	Upload(userID uuid.UUID, file *multipart.FileHeader) (*models.CVDocument, error)
	Get(userID uuid.UUID) (*models.CVDocument, error)
	Delete(userID uuid.UUID) error
}

type cvService struct {
	repo           repository.CVRepository
	profileService ProfileService
	config         *config.Config
}

func NewCVService(repo repository.CVRepository, profileService ProfileService, cfg *config.Config) CVService {
	return &cvService{
		repo:           repo,
		profileService: profileService,
		config:         cfg,
	}
}

func (s *cvService) Upload(userID uuid.UUID, file *multipart.FileHeader) (*models.CVDocument, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	// Validate file size
	if file.Size > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxFileSize)
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(s.config.AllowedFileTypes, ext) {
		return nil, errors.New("file type not allowed. Allowed types: " + s.config.AllowedFileTypes)
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(s.config.UploadDir, 0755); err != nil {
		return nil, err
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", profile.ID.String(), time.Now().Unix(), ext)
	filePath := filepath.Join(s.config.UploadDir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Check if CV already exists for this profile
	existing, _ := s.repo.GetByProfileID(profile.ID)

	cvDoc := &models.CVDocument{
		ProfileID:  profile.ID,
		FileName:   file.Filename,
		FileURL:    filePath,
		FileSize:   file.Size,
		FileType:   ext,
		IsVerified: false,
		UploadedAt: time.Now(),
	}

	if existing != nil && existing.ID != uuid.Nil {
		// Update existing CV
		cvDoc.ID = existing.ID

		// Delete old file
		if existing.FileURL != "" {
			os.Remove(existing.FileURL)
		}

		err = s.repo.Update(cvDoc)
	} else {
		// Create new CV record
		err = s.repo.Create(cvDoc)
	}

	if err != nil {
		// Cleanup uploaded file if database operation fails
		os.Remove(filePath)
		return nil, err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return cvDoc, nil
}

func (s *cvService) Get(userID uuid.UUID) (*models.CVDocument, error) {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	cv, err := s.repo.GetByProfileID(profile.ID)
	if err != nil {
		return nil, errors.New("CV not found")
	}
	return cv, nil
}

func (s *cvService) Delete(userID uuid.UUID) error {
	profile, err := s.profileService.GetProfile(userID)
	if err != nil {
		return errors.New("profile not found")
	}

	cv, err := s.repo.GetByProfileID(profile.ID)
	if err != nil {
		return errors.New("CV not found")
	}

	// Delete file from storage
	if cv.FileURL != "" {
		os.Remove(cv.FileURL)
	}

	err = s.repo.Delete(cv.ID)
	if err != nil {
		return err
	}

	s.profileService.UpdateCompletionStatus(userID)
	return nil
}
