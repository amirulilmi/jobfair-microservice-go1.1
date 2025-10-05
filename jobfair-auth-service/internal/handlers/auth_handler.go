package handlers

import (
	"net/http"

	"jobfair-auth-service/internal/models"
	"jobfair-auth-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user": gin.H{
				"id":        user.ID,
				"email":     user.Email,
				"user_type": user.UserType,
				"is_active": user.IsActive,
			},
		},
		"message": "User registered successfully",
		"status":  true,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":         response.Token,
			"refresh_token": response.RefreshToken,
			"user": gin.H{
				"id":        response.User.ID,
				"email":     response.User.Email,
				"user_type": response.User.UserType,
				"is_active": response.User.IsActive,
			},
		},
		"message": "Login successful",
		"status":  true,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Refresh token required in header",
		})
		return
	}

	token, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token": token,
		},
		"message": "Token refreshed successfully",
		"status":  true,
	})
}

func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetCurrentUser returns the complete profile of the logged-in user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user_id from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized - user_id not found in context",
		})
		return
	}

	// Get complete user data with profile
	userData, err := h.authService.GetCurrentUserWithProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User data retrieved successfully",
		"data":    userData,
	})
}
