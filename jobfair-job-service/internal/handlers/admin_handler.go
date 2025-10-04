package handlers

import (
	"fmt"
	"log"
	"net/http"

	"jobfair-job-service/internal/repository"
	"jobfair-job-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	companyRepo *repository.CompanyRepository
	jobService  *services.JobService
}

func NewAdminHandler(
	companyRepo *repository.CompanyRepository,
	jobService *services.JobService,
) *AdminHandler {
	return &AdminHandler{
		companyRepo: companyRepo,
		jobService:  jobService,
	}
}

// SyncCompanyMapping manually syncs a company mapping
// POST /api/v1/admin/sync-company-mapping
func (h *AdminHandler) SyncCompanyMapping(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id" binding:"required"`
		CompanyID   uint   `json:"company_id" binding:"required"`
		CompanyName string `json:"company_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Validate - only admin can do this
	userType := c.GetString("user_type")
	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Only admins can sync company mappings",
		})
		return
	}

	// Create or update mapping
	if err := h.companyRepo.UpsertCompanyMapping(req.UserID, req.CompanyID, req.CompanyName); err != nil {
		log.Printf("❌ Failed to sync company mapping: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to sync company mapping",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("✅ [ADMIN] Company mapping synced: UserID=%d, CompanyID=%d, Name=%s",
		req.UserID, req.CompanyID, req.CompanyName)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Company mapping synced successfully",
		"data": gin.H{
			"user_id":      req.UserID,
			"company_id":   req.CompanyID,
			"company_name": req.CompanyName,
		},
	})
}

// HealthCheckDataConsistency checks for data inconsistencies
// GET /api/v1/admin/health/data-consistency
func (h *AdminHandler) HealthCheckDataConsistency(c *gin.Context) {
	// Validate - only admin can do this
	userType := c.GetString("user_type")
	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Only admins can check data consistency",
		})
		return
	}

	issues := []string{}
	warnings := []string{}

	// Check 1: Jobs without valid company mappings
	var jobsWithoutMapping int64
	h.companyRepo.GetDB().Raw(`
		SELECT COUNT(*) 
		FROM jobs j
		LEFT JOIN company_mappings cm ON j.user_id = cm.user_id
		WHERE cm.id IS NULL AND j.deleted_at IS NULL
	`).Scan(&jobsWithoutMapping)

	if jobsWithoutMapping > 0 {
		issues = append(issues, fmt.Sprintf("%d jobs without valid company mapping", jobsWithoutMapping))
	}

	// Check 2: Orphaned applications
	var orphanedApplications int64
	h.companyRepo.GetDB().Raw(`
		SELECT COUNT(*)
		FROM job_applications ja
		LEFT JOIN jobs j ON ja.job_id = j.id
		WHERE j.id IS NULL AND ja.deleted_at IS NULL
	`).Scan(&orphanedApplications)

	if orphanedApplications > 0 {
		warnings = append(warnings, fmt.Sprintf("%d orphaned applications", orphanedApplications))
	}

	// Check 3: Jobs with deleted companies
	var jobsWithDeletedCompany int64
	h.companyRepo.GetDB().Raw(`
		SELECT COUNT(*)
		FROM jobs j
		INNER JOIN company_mappings cm ON j.user_id = cm.user_id
		WHERE j.deleted_at IS NULL
		-- Note: Can't check company deletion without dblink to company-service DB
	`).Scan(&jobsWithDeletedCompany)

	status := "healthy"
	if len(issues) > 0 {
		status = "unhealthy"
	} else if len(warnings) > 0 {
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  status,
		"checks": gin.H{
			"jobs_without_mapping":  jobsWithoutMapping,
			"orphaned_applications": orphanedApplications,
		},
		"issues":   issues,
		"warnings": warnings,
	})
}

// GetCompanyMappings lists all company mappings (admin only)
// GET /api/v1/admin/company-mappings
func (h *AdminHandler) GetCompanyMappings(c *gin.Context) {
	// Validate - only admin
	userType := c.GetString("user_type")
	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Only admins can view all company mappings",
		})
		return
	}

	var mappings []repository.CompanyMapping
	if err := h.companyRepo.GetDB().Find(&mappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve company mappings",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mappings,
		"count":   len(mappings),
	})
}
