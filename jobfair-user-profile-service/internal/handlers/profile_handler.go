package handlers

import (
	"net/http"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	service services.ProfileService
}

func NewProfileHandler(service services.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req struct {
		FullName    string `json:"full_name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	profile, err := h.service.CreateProfile(userID.(uint), req.FullName, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Profile created successfully", profile))
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	profile, err := h.service.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Profile not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Profile retrieved successfully", profile))
}

func (h *ProfileHandler) GetProfileWithRelations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	profile, err := h.service.GetProfileWithRelations(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Profile not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Profile retrieved successfully", profile))
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	profile, err := h.service.UpdateProfile(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Profile updated successfully", profile))
}

func (h *ProfileHandler) GetCompletionStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	profile, err := h.service.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Profile not found", "NOT_FOUND", nil))
		return
	}

	completionStatus := h.service.CalculateCompletionStatus(profile)

	c.JSON(http.StatusOK, models.SuccessResponse("Completion status retrieved", gin.H{
		"completion_status": completionStatus,
		"is_complete":       completionStatus == 100,
	}))
}
