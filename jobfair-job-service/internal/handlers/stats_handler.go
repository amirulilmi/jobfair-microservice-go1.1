package handlers

import (
	"net/http"
	"strconv"

	"jobfair-job-service/internal/models"
	"jobfair-job-service/internal/services"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	applicationService *services.ApplicationService
	jobService         *services.JobService
}

func NewStatsHandler(applicationService *services.ApplicationService, jobService *services.JobService) *StatsHandler {
	return &StatsHandler{
		applicationService: applicationService,
		jobService:         jobService,
	}
}

// GetUserApplicationsCount returns count of applications by user
func (h *StatsHandler) GetUserApplicationsCount(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	count, err := h.applicationService.GetUserApplicationsCount(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    count,
	})
}

// GetUserSavedCount returns count of saved jobs by user
func (h *StatsHandler) GetUserSavedCount(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	count, err := h.jobService.GetUserSavedCount(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    count,
	})
}
