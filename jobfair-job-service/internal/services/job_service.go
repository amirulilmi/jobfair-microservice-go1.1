package services

import (
	"errors"
	"fmt"
	"jobfair-job-service/internal/models"
	"jobfair-job-service/internal/repository"
	"jobfair-job-service/internal/utils"
	"time"

	"gorm.io/gorm"
)

type JobService struct {
	jobRepo         *repository.JobRepository
	applicationRepo *repository.ApplicationRepository
	savedJobRepo    *repository.SavedJobRepository
}

func NewJobService(
	jobRepo *repository.JobRepository,
	applicationRepo *repository.ApplicationRepository,
	savedJobRepo *repository.SavedJobRepository,
) *JobService {
	return &JobService{
		jobRepo:         jobRepo,
		applicationRepo: applicationRepo,
		savedJobRepo:    savedJobRepo,
	}
}

// CreateJob creates a new job posting
func (s *JobService) CreateJob(userID, companyID uint, req *models.CreateJobRequest) (*models.Job, error) {
	// Generate slug from title
	slug := utils.GenerateSlug(req.Title)

	job := &models.Job{
		CompanyID:        companyID,
		UserID:           userID,
		Title:            req.Title,
		Description:      req.Description,
		Slug:             slug,
		EmploymentType:   req.EmploymentType,
		WorkType:         req.WorkType,
		ExperienceLevel:  req.ExperienceLevel,
		Location:         req.Location,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		Skills:           req.Skills,
		Benefits:         req.Benefits,
		ReceiveMethod:    req.ReceiveMethod,
		ContactEmail:     req.ContactEmail,
		ExternalURL:      req.ExternalURL,
		Tags:             req.Tags,
		Status:           models.JobStatusDraft,
	}

	// Parse deadline
	if req.Deadline != nil && *req.Deadline != "" {
		deadline, err := time.Parse(time.RFC3339, *req.Deadline)
		if err == nil {
			job.Deadline = &deadline
		}
	}

	return s.jobRepo.Create(job)
}

// GetJobByID retrieves a job by ID with additional context
func (s *JobService) GetJobByID(jobID uint, userID *uint) (*models.JobDetailResponse, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	response := &models.JobDetailResponse{
		Job: job,
	}

	// Check if user has saved or applied
	if userID != nil {
		isSaved, _ := s.savedJobRepo.IsSaved(jobID, *userID)
		response.IsSaved = isSaved

		application, err := s.applicationRepo.GetByJobIDAndUserID(jobID, *userID)
		if err == nil {
			response.HasApplied = true
			response.Application = application
		}
	}

	// Increment views
	go s.jobRepo.IncrementViews(jobID)

	return response, nil
}

// GetJobBySlug retrieves a job by slug
func (s *JobService) GetJobBySlug(slug string, userID *uint) (*models.JobDetailResponse, error) {
	job, err := s.jobRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.GetJobByID(job.ID, userID)
}

// UpdateJob updates a job
func (s *JobService) UpdateJob(jobID, userID uint, req *models.UpdateJobRequest) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if job.UserID != userID {
		return nil, errors.New("unauthorized to update this job")
	}

	// Update fields
	if req.Title != nil {
		job.Title = *req.Title
		job.Slug = utils.GenerateSlug(*req.Title)
	}
	if req.Description != nil {
		job.Description = *req.Description
	}
	if req.EmploymentType != nil {
		job.EmploymentType = *req.EmploymentType
	}
	if req.WorkType != nil {
		job.WorkType = *req.WorkType
	}
	if req.ExperienceLevel != nil {
		job.ExperienceLevel = *req.ExperienceLevel
	}
	if req.Location != nil {
		job.Location = *req.Location
	}
	if req.SalaryMin != nil {
		job.SalaryMin = *req.SalaryMin
	}
	if req.SalaryMax != nil {
		job.SalaryMax = *req.SalaryMax
	}
	if req.Requirements != nil {
		job.Requirements = req.Requirements
	}
	if req.Responsibilities != nil {
		job.Responsibilities = req.Responsibilities
	}
	if req.Skills != nil {
		job.Skills = req.Skills
	}
	if req.Benefits != nil {
		job.Benefits = req.Benefits
	}
	if req.ReceiveMethod != nil {
		job.ReceiveMethod = *req.ReceiveMethod
	}
	if req.ContactEmail != nil {
		job.ContactEmail = *req.ContactEmail
	}
	if req.ExternalURL != nil {
		job.ExternalURL = req.ExternalURL
	}
	if req.Tags != nil {
		job.Tags = req.Tags
	}
	if req.Deadline != nil && *req.Deadline != "" {
		deadline, err := time.Parse(time.RFC3339, *req.Deadline)
		if err == nil {
			job.Deadline = &deadline
		}
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// DeleteJob soft deletes a job
func (s *JobService) DeleteJob(jobID, userID uint) error {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return err
	}

	// Check ownership
	if job.UserID != userID {
		return errors.New("unauthorized to delete this job")
	}

	return s.jobRepo.Delete(jobID)
}

// PublishJob publishes a draft job
func (s *JobService) PublishJob(jobID, userID uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if job.UserID != userID {
		return nil, errors.New("unauthorized to publish this job")
	}

	if job.Status != models.JobStatusDraft {
		return nil, fmt.Errorf("job is already %s", job.Status)
	}

	now := time.Now()
	job.Status = models.JobStatusPublished
	job.PublishedAt = &now

	if err := s.jobRepo.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// CloseJob closes a job
func (s *JobService) CloseJob(jobID, userID uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if job.UserID != userID {
		return nil, errors.New("unauthorized to close this job")
	}

	if job.Status == models.JobStatusClosed {
		return nil, errors.New("job is already closed")
	}

	now := time.Now()
	job.Status = models.JobStatusClosed
	job.ClosedAt = &now

	if err := s.jobRepo.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

// ListJobs retrieves jobs with filters
func (s *JobService) ListJobs(filter models.JobListFilter) ([]*models.Job, *models.PaginationMeta, error) {
	// Default to published jobs for public listing
	if filter.Status == "" {
		filter.Status = models.JobStatusPublished
	}

	jobs, total, err := s.jobRepo.List(filter)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: filter.Page,
		PerPage:     filter.Limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return jobs, meta, nil
}

// GetMyJobs retrieves jobs posted by the user
func (s *JobService) GetMyJobs(userID uint) ([]*models.Job, error) {
	return s.jobRepo.GetByUserID(userID)
}

// ApplyToJob applies to a job
func (s *JobService) ApplyToJob(jobID, userID uint, req *models.ApplyJobRequest) (*models.JobApplication, error) {
	// Check if job exists
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	// Check if job is accepting applications
	if job.Status != models.JobStatusPublished {
		return nil, errors.New("job is not accepting applications")
	}

	// Check if already applied
	_, err = s.applicationRepo.GetByJobIDAndUserID(jobID, userID)
	if err == nil {
		return nil, errors.New("already applied to this job")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create application
	application := &models.JobApplication{
		JobID:       jobID,
		UserID:      userID,
		CVURL:       req.CVURL,
		CoverLetter: req.CoverLetter,
		Status:      models.ApplicationStatusApplied,
	}

	createdApp, err := s.applicationRepo.Create(application)
	if err != nil {
		return nil, err
	}

	// Increment applications count
	go s.jobRepo.IncrementApplications(jobID)

	return createdApp, nil
}

// BulkApply applies to multiple jobs at once
func (s *JobService) BulkApply(userID uint, req *models.BulkApplyRequest) ([]uint, []error) {
	var successIDs []uint
	var errors []error

	for _, jobID := range req.JobIDs {
		applyReq := &models.ApplyJobRequest{
			CVURL:       req.CVURL,
			CoverLetter: req.CoverLetter,
		}

		_, err := s.ApplyToJob(jobID, userID, applyReq)
		if err != nil {
			errors = append(errors, fmt.Errorf("job %d: %w", jobID, err))
		} else {
			successIDs = append(successIDs, jobID)
		}
	}

	return successIDs, errors
}

// SaveJob bookmarks a job
func (s *JobService) SaveJob(jobID, userID uint) error {
	// Check if job exists
	_, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return err
	}

	// Check if already saved
	_, err = s.savedJobRepo.GetByJobIDAndUserID(jobID, userID)
	if err == nil {
		return errors.New("job already saved")
	}

	savedJob := &models.SavedJob{
		JobID:  jobID,
		UserID: userID,
	}

	_, err = s.savedJobRepo.Create(savedJob)
	return err
}

// UnsaveJob removes bookmark
func (s *JobService) UnsaveJob(jobID, userID uint) error {
	return s.savedJobRepo.DeleteByJobIDAndUserID(jobID, userID)
}

// GetSavedJobs retrieves saved jobs
func (s *JobService) GetSavedJobs(userID uint, page, limit int) ([]*models.SavedJob, *models.PaginationMeta, error) {
	savedJobs, total, err := s.savedJobRepo.GetByUserID(userID, page, limit)
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

	return savedJobs, meta, nil
}
