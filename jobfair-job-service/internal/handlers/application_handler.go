package handlers

import (
	"net/http"
	"strconv"

	"jobfair-job-service/internal/models"
	"jobfair-job-service/internal/services"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	applicationService *services.ApplicationService
}

func NewApplicationHandler(applicationService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService}
}

// GetApplication handles GET /applications/:id
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid application ID",
		})
		return
	}

	isCompany := userType == "company"
	application, err := h.applicationService.GetApplicationByID(uint(id), userID, isCompany)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    application,
	})
}

// GetMyApplications handles GET /applications/my
func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "job_seeker" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only job seekers can view their applications",
		})
		return
	}

	status := models.ApplicationStatus(c.Query("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	applications, meta, err := h.applicationService.GetMyApplications(userID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    applications,
		Meta:    meta,
	})
}

// GetApplicationsByJobID handles GET /jobs/:id/applications
func (h *ApplicationHandler) GetApplicationsByJobID(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "company" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only companies can view job applications",
		})
		return
	}

	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid job ID",
		})
		return
	}

	status := models.ApplicationStatus(c.Query("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	applications, meta, err := h.applicationService.GetApplicationsByJobID(uint(jobID), userID, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    applications,
		Meta:    meta,
	})
}

// UpdateApplicationStatus handles PUT /applications/:id/status
func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "company" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only companies can update application status",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid application ID",
		})
		return
	}

	var req models.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	application, err := h.applicationService.UpdateApplicationStatus(uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Application status updated successfully",
		Data:    application,
	})
}

// WithdrawApplication handles DELETE /applications/:id
func (h *ApplicationHandler) WithdrawApplication(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType != "job_seeker" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only job seekers can withdraw their applications",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid application ID",
		})
		return
	}

	if err := h.applicationService.WithdrawApplication(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Application withdrawn successfully",
	})
}

// GetApplicationStats handles GET /applications/stats
func (h *ApplicationHandler) GetApplicationStats(c *gin.Context) {
	userType := c.GetString("user_type")

	if userType != "company" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only companies can view application stats",
		})
		return
	}

	companyIDStr := c.Query("company_id")
	if companyIDStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "company_id is required",
		})
		return
	}

	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid company ID",
		})
		return
	}

	stats, err := h.applicationService.GetApplicationStats(uint(companyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}
