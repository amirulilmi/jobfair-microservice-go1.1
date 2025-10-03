package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
)

type PreferenceHandler struct {
	service services.PreferenceService
}

func NewPreferenceHandler(service services.PreferenceService) *PreferenceHandler {
	return &PreferenceHandler{service: service}
}

func (h *PreferenceHandler) CreateOrUpdateCareerPreference(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.CareerPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	// Convert arrays to comma-separated strings
	preferredWorkTypes := strings.Join(req.PreferredWorkTypes, ",")
	preferredLocations := strings.Join(req.PreferredLocations, ",")

	preference, err := h.service.CreateOrUpdateCareerPreference(
		userID.(uint),
		req.IsActivelyLooking,
		req.ExpectedSalaryMin,
		req.ExpectedSalaryMax,
		req.SalaryCurrency,
		req.IsNegotiable,
		preferredWorkTypes,
		preferredLocations,
		req.WillingToRelocate,
		req.AvailableStartDate,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "SAVE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Career preference saved successfully", preference))
}

func (h *PreferenceHandler) GetCareerPreference(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	preference, err := h.service.GetCareerPreference(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Career preference not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Career preference retrieved successfully", preference))
}

func (h *PreferenceHandler) CreatePositionPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.BulkPositionPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	preferences, err := h.service.CreatePositionPreferences(userID.(uint), req.Positions)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Position preferences created successfully", preferences))
}

func (h *PreferenceHandler) GetPositionPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	preferences, err := h.service.GetPositionPreferences(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "FETCH_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Position preferences retrieved successfully", preferences))
}

func (h *PreferenceHandler) DeletePositionPreference(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid ID", "INVALID_ID", nil))
		return
	}

	err = h.service.DeletePositionPreference(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "DELETE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Position preference deleted successfully", nil))
}
