package handlers

import (
	"net/http"
	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SkillHandler struct {
	service services.SkillService
}

func NewSkillHandler(service services.SkillService) *SkillHandler {
	return &SkillHandler{service: service}
}

func (h *SkillHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.SkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	skill, err := h.service.Create(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Skill created successfully", skill))
}

func (h *SkillHandler) CreateBulk(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	var req models.BulkSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	skills, err := h.service.CreateBulk(userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Skills created successfully", skills))
}

func (h *SkillHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	skills, err := h.service.GetAll(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error(), "FETCH_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Skills retrieved successfully", skills))
}

func (h *SkillHandler) GetByID(c *gin.Context) {
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

	skill, err := h.service.GetByID(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Skill not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Skill retrieved successfully", skill))
}

func (h *SkillHandler) Update(c *gin.Context) {
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

	var req models.SkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	skill, err := h.service.Update(userID.(uuid.UUID), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Skill updated successfully", skill))
}

func (h *SkillHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, models.SuccessResponse("Skill deleted successfully", nil))
}
