package services

import (
	"errors"
	"jobfair-job-service/internal/models"
	"jobfair-job-service/internal/repository"
	"time"
)

type ApplicationService struct {
	applicationRepo *repository.ApplicationRepository
	jobRepo         *repository.JobRepository
}

func NewApplicationService(
	applicationRepo *repository.ApplicationRepository,
	jobRepo *repository.JobRepository,
) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
	}
}

// GetApplicationByID retrieves an application by ID
func (s *ApplicationService) GetApplicationByID(applicationID, userID uint, isCompany bool) (*models.JobApplication, error) {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return nil, err
	}

	// Check authorization
	if !isCompany && application.UserID != userID {
		return nil, errors.New("unauthorized to view this application")
	}

	if isCompany {
		job, err := s.jobRepo.GetByID(application.JobID)
		if err != nil {
			return nil, err
		}
		if job.UserID != userID {
			return nil, errors.New("unauthorized to view this application")
		}
	}

	return application, nil
}

// GetMyApplications retrieves all applications by a user
func (s *ApplicationService) GetMyApplications(userID uint, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, *models.PaginationMeta, error) {
	applications, total, err := s.applicationRepo.GetByUserID(userID, status, page, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return applications, meta, nil
}

// GetApplicationsByJobID retrieves applications for a specific job
func (s *ApplicationService) GetApplicationsByJobID(jobID, userID uint, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, *models.PaginationMeta, error) {
	// Check if user owns the job
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, nil, err
	}

	if job.UserID != userID {
		return nil, nil, errors.New("unauthorized to view applications for this job")
	}

	applications, total, err := s.applicationRepo.GetByJobID(jobID, status, page, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return applications, meta, nil
}

// GetApplicationsByCompanyID retrieves all applications for a company's jobs
func (s *ApplicationService) GetApplicationsByCompanyID(companyID uint, status models.ApplicationStatus, page, limit int) ([]*models.JobApplication, *models.PaginationMeta, error) {
	applications, total, err := s.applicationRepo.GetByCompanyID(companyID, status, page, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return applications, meta, nil
}

// UpdateApplicationStatus updates the status of an application
func (s *ApplicationService) UpdateApplicationStatus(applicationID, userID uint, req *models.UpdateApplicationStatusRequest) (*models.JobApplication, error) {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the job
	job, err := s.jobRepo.GetByID(application.JobID)
	if err != nil {
		return nil, err
	}

	if job.UserID != userID {
		return nil, errors.New("unauthorized to update this application")
	}

	// Update status
	application.Status = req.Status
	application.StatusNote = req.StatusNote

	// Update tracking timestamps
	now := time.Now()
	switch req.Status {
	case models.ApplicationStatusReviewing:
		if application.ViewedAt == nil {
			application.ViewedAt = &now
		}
		application.ReviewedAt = &now
	case models.ApplicationStatusInterview:
		application.InterviewAt = &now
	case models.ApplicationStatusHired, models.ApplicationStatusRejected:
		application.RespondedAt = &now
	}

	if err := s.applicationRepo.Update(application); err != nil {
		return nil, err
	}

	return application, nil
}

// WithdrawApplication withdraws an application
func (s *ApplicationService) WithdrawApplication(applicationID, userID uint) error {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return err
	}

	// Check ownership
	if application.UserID != userID {
		return errors.New("unauthorized to withdraw this application")
	}

	// Delete application
	if err := s.applicationRepo.Delete(applicationID); err != nil {
		return err
	}

	// Decrement job applications count
	go s.jobRepo.DecrementApplications(application.JobID)

	return nil
}

// GetApplicationStats retrieves application statistics for a company
func (s *ApplicationService) GetApplicationStats(companyID uint) (map[string]int64, error) {
	return s.applicationRepo.GetApplicationStats(companyID)
}

// GetUserApplicationsCount returns count of applications by user
func (s *ApplicationService) GetUserApplicationsCount(userID uint) (int64, error) {
	return s.applicationRepo.CountByUserID(userID)
}
