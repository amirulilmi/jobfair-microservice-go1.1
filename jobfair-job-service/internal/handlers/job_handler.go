package handlers

import (
	"net/http"
	"strconv"

	"jobfair-job-service/internal/models"
	"jobfair-job-service/internal/services"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

// CreateJob handles POST /jobs
func (h *JobHandler) CreateJob(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "company" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only companies can post jobs",
		})
		return
	}

	// Auto-detect company_id from user_id
	companyID, err := h.jobService.GetCompanyIDByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Company profile not found. Please complete company registration first.",
		})
		return
	}

	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	job, err := h.jobService.CreateJob(userID, companyID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Job created successfully",
		Data:    job,
	})
}

// GetJob handles GET /jobs/:id
func (h *JobHandler) GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	// Get user ID if authenticated
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		uidValue := uid.(uint)
		userID = &uidValue
	}

	jobDetail, err := h.jobService.GetJobByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Job not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    jobDetail,
	})
}

// UpdateJob handles PUT /jobs/:id
func (h *JobHandler) UpdateJob(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	var req models.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	job, err := h.jobService.UpdateJob(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Job updated successfully",
		Data:    job,
	})
}

// DeleteJob handles DELETE /jobs/:id
func (h *JobHandler) DeleteJob(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	if err := h.jobService.DeleteJob(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Job deleted successfully",
	})
}

// PublishJob handles POST /jobs/:id/publish
func (h *JobHandler) PublishJob(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	job, err := h.jobService.PublishJob(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Job published successfully",
		Data:    job,
	})
}

// CloseJob handles POST /jobs/:id/close
func (h *JobHandler) CloseJob(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	job, err := h.jobService.CloseJob(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Job closed successfully",
		Data:    job,
	})
}

// ListJobs handles GET /jobs
func (h *JobHandler) ListJobs(c *gin.Context) {
	var filter models.JobListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Default pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	jobs, meta, err := h.jobService.ListJobs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Enrich jobs with company data
	jobsWithCompany, err := h.jobService.EnrichJobsWithCompanyData(jobs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to enrich job data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    jobsWithCompany,
		Meta:    meta,
	})
}

// GetMyJobs handles GET /jobs/my
func (h *JobHandler) GetMyJobs(c *gin.Context) {
	userID := c.GetUint("user_id")

	jobs, err := h.jobService.GetMyJobs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    jobs,
	})
}

// ApplyToJob handles POST /jobs/:id/apply
func (h *JobHandler) ApplyToJob(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "job_seeker" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only job seekers can apply to jobs",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	var req models.ApplyJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	application, err := h.jobService.ApplyToJob(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Application submitted successfully",
		Data:    application,
	})
}

// BulkApply handles POST /jobs/bulk-apply
func (h *JobHandler) BulkApply(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "job_seeker" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only job seekers can apply to jobs",
		})
		return
	}

	var req models.BulkApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	successIDs, errors := h.jobService.BulkApply(userID, &req)

	response := map[string]interface{}{
		"success_count": len(successIDs),
		"success_ids":   successIDs,
		"error_count":   len(errors),
	}

	if len(errors) > 0 {
		errorMessages := make([]string, len(errors))
		for i, err := range errors {
			errorMessages[i] = err.Error()
		}
		response["errors"] = errorMessages
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bulk application completed",
		Data:    response,
	})
}

// SaveJob handles POST /jobs/:id/save
func (h *JobHandler) SaveJob(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	if err := h.jobService.SaveJob(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Job saved successfully",
	})
}

// UnsaveJob handles DELETE /jobs/:id/save
func (h *JobHandler) UnsaveJob(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	if err := h.jobService.UnsaveJob(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Job unsaved successfully",
	})
}

// GetSavedJobs handles GET /jobs/saved
func (h *JobHandler) GetSavedJobs(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	savedJobs, meta, err := h.jobService.GetSavedJobs(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    savedJobs,
		Meta:    meta,
	})
}
