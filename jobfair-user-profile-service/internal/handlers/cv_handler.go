package handlers

import (
	"net/http"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CVHandler struct {
	service services.CVService
}

func NewCVHandler(service services.CVService) *CVHandler {
	return &CVHandler{service: service}
}

func (h *CVHandler) Upload(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("No file uploaded", "NO_FILE", nil))
		return
	}

	cv, err := h.service.Upload(userID.(uuid.UUID), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPLOAD_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("CV uploaded successfully", cv))
}

func (h *CVHandler) Get(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	cv, err := h.service.Get(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("CV not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("CV retrieved successfully", cv))
}

func (h *CVHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	err := h.service.Delete(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "DELETE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("CV deleted successfully", nil))
}
