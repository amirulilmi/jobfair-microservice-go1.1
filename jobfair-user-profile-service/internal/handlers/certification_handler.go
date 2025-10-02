package handlers

import (
	"net/http"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CertificationHandler struct {
	service services.CertificationService
}

func NewCertificationHandler(service services.CertificationService) *CertificationHandler {
	return &CertificationHandler{service: service}
}

func (h *CertificationHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.CertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	certification, err := h.service.Create(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Certification created successfully", certification))
}

func (h *CertificationHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	certifications, err := h.service.GetAll(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "FETCH_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Certifications retrieved successfully", certifications))
}

func (h *CertificationHandler) GetByID(c *gin.Context) {
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

	certification, err := h.service.GetByID(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Certification not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Certification retrieved successfully", certification))
}

func (h *CertificationHandler) Update(c *gin.Context) {
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

	var req models.CertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	certification, err := h.service.Update(userID.(uuid.UUID), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Certification updated successfully", certification))
}

func (h *CertificationHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, models.SuccessResponse("Certification deleted successfully", nil))
}
