package handlers

import (
	"net/http"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	service services.ProfileService
}

func NewBannerHandler(service services.ProfileService) *BannerHandler {
	return &BannerHandler{service: service}
}

func (h *BannerHandler) UploadBanner(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	file, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("No banner file uploaded", "NO_FILE", nil))
		return
	}

	bannerURL, err := h.service.UploadBanner(userID.(uint), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPLOAD_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Banner uploaded successfully", gin.H{
		"banner_image_url": bannerURL,
	}))
}

func (h *BannerHandler) DeleteBanner(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	err := h.service.DeleteBanner(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "DELETE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Banner deleted successfully", nil))
}
