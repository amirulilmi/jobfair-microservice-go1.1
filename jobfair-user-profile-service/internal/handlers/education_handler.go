package handlers

import (
	"net/http"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EducationHandler struct {
	service services.EducationService
}

func NewEducationHandler(service services.EducationService) *EducationHandler {
	return &EducationHandler{service: service}
}

func (h *EducationHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.EducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	education, err := h.service.Create(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Education created successfully", education))
}

func (h *EducationHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	educations, err := h.service.GetAll(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "FETCH_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Educations retrieved successfully", educations))
}

func (h *EducationHandler) GetByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid ID", "INVALID_ID", nil))
		return
	}

	education, err := h.service.GetByID(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Education not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Education retrieved successfully", education))
}

func (h *EducationHandler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid ID", "INVALID_ID", nil))
		return
	}

	var req models.EducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	education, err := h.service.Update(userID.(uuid.UUID), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Education updated successfully", education))
}

func (h *EducationHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid ID", "INVALID_ID", nil))
		return
	}

	err = h.service.Delete(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "DELETE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Education deleted successfully", nil))
}
