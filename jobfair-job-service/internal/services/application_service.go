package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jobfair-job-service/internal/models"
	"jobfair-job-service/internal/repository"
	"net/http"
	"time"
)

type ApplicationService struct {
	applicationRepo   *repository.ApplicationRepository
	jobRepo           *repository.JobRepository
	companyServiceURL string
}

func NewApplicationService(
	applicationRepo *repository.ApplicationRepository,
	jobRepo *repository.JobRepository, companyServiceURL string,
) *ApplicationService {
	return &ApplicationService{
		applicationRepo:   applicationRepo,
		jobRepo:           jobRepo,
		companyServiceURL: companyServiceURL,
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

// EnrichApplicationsWithCompanyData enriches applications with company data
func (s *ApplicationService) EnrichApplicationsWithCompanyData(applications []*models.JobApplication) ([]models.ApplicationWithCompany, error) {
	if len(applications) == 0 {
		return []models.ApplicationWithCompany{}, nil
	}

	// Collect unique company IDs from jobs
	companyIDs := make(map[uint]bool)
	for _, app := range applications {
		if app.Job != nil {
			companyIDs[app.Job.CompanyID] = true
		}
	}

	// Fetch company data for all unique IDs
	companyDataMap := make(map[uint]map[string]interface{})
	for companyID := range companyIDs {
		companyData, err := s.fetchCompanyData(companyID)
		if err != nil {
			fmt.Printf("[WARNING] Failed to fetch company %d: %v\n", companyID, err)
			// Use fallback company data if fetch fails
			companyDataMap[companyID] = map[string]interface{}{
				"id":   companyID,
				"name": "Unknown Company",
			}
		} else {
			companyDataMap[companyID] = companyData
		}
	}

	// Enrich applications with company data
	result := make([]models.ApplicationWithCompany, len(applications))
	for i, app := range applications {
		var companyData map[string]interface{}
		if app.Job != nil {
			companyData = companyDataMap[app.Job.CompanyID]
		}

		result[i] = models.ApplicationWithCompany{
			JobApplication: app,
			Job:            app.Job,
			Company:        companyData,
		}
	}

	return result, nil
}

// fetchCompanyData fetches company data from company service
func (s *ApplicationService) fetchCompanyData(companyID uint) (map[string]interface{}, error) {
	// Get company service URL from job service
	// We need to add this as a dependency
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/api/v1/companies/%d", s.companyServiceURL, companyID)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("company service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New("company service returned unsuccessful response")
	}

	return result.Data, nil
}
