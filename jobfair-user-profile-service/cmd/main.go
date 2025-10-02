package main

import (
	"jobfair-user-profile-service/internal/config"
	"jobfair-user-profile-service/internal/handlers"
	"jobfair-user-profile-service/internal/middleware"
	"jobfair-user-profile-service/internal/repository"
	"jobfair-user-profile-service/internal/services"
	"jobfair-user-profile-service/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// Initialize repositories
	profileRepo := repository.NewProfileRepository(db)
	workExpRepo := repository.NewWorkExperienceRepository(db)
	educationRepo := repository.NewEducationRepository(db)
	certificationRepo := repository.NewCertificationRepository(db)
	skillRepo := repository.NewSkillRepository(db)
	preferenceRepo := repository.NewPreferenceRepository(db)
	cvRepo := repository.NewCVRepository(db)
	// badgeRepo := repository.NewBadgeRepository(db)

	// Initialize services
	profileService := services.NewProfileService(
		profileRepo,
		workExpRepo,
		educationRepo,
		certificationRepo,
		skillRepo,
		preferenceRepo,
		cvRepo,
	)
	workExpService := services.NewWorkExperienceService(workExpRepo, profileService)
	educationService := services.NewEducationService(educationRepo, profileService)
	certificationService := services.NewCertificationService(certificationRepo, profileService)
	skillService := services.NewSkillService(skillRepo, profileService)
	preferenceService := services.NewPreferenceService(preferenceRepo, profileService)
	cvService := services.NewCVService(cvRepo, profileService, cfg)

	// Initialize handlers
	profileHandler := handlers.NewProfileHandler(profileService)
	workExpHandler := handlers.NewWorkExperienceHandler(workExpService)
	educationHandler := handlers.NewEducationHandler(educationService)
	certificationHandler := handlers.NewCertificationHandler(certificationService)
	skillHandler := handlers.NewSkillHandler(skillService)
	preferenceHandler := handlers.NewPreferenceHandler(preferenceService)
	cvHandler := handlers.NewCVHandler(cvService)

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": cfg.ServiceName,
		})
	})

	// API v1 routes with JWT authentication
	v1 := router.Group("/api/v1")
	v1.Use(middleware.JWTAuthMiddleware(cfg.JWTSecret))
	{
		// Profile routes
		profiles := v1.Group("/profiles")
		{
			profiles.POST("", profileHandler.CreateProfile)
			profiles.GET("", profileHandler.GetProfile)
			profiles.GET("/full", profileHandler.GetProfileWithRelations)
			profiles.PUT("", profileHandler.UpdateProfile)
			profiles.GET("/completion", profileHandler.GetCompletionStatus)
		}

		// Work Experience routes
		workExps := v1.Group("/work-experiences")
		{
			workExps.POST("", workExpHandler.Create)
			workExps.GET("", workExpHandler.GetAll)
			workExps.GET("/:id", workExpHandler.GetByID)
			workExps.PUT("/:id", workExpHandler.Update)
			workExps.DELETE("/:id", workExpHandler.Delete)
		}

		// Education routes
		educations := v1.Group("/educations")
		{
			educations.POST("", educationHandler.Create)
			educations.GET("", educationHandler.GetAll)
			educations.GET("/:id", educationHandler.GetByID)
			educations.PUT("/:id", educationHandler.Update)
			educations.DELETE("/:id", educationHandler.Delete)
		}

		// Certification routes
		certifications := v1.Group("/certifications")
		{
			certifications.POST("", certificationHandler.Create)
			certifications.GET("", certificationHandler.GetAll)
			certifications.GET("/:id", certificationHandler.GetByID)
			certifications.PUT("/:id", certificationHandler.Update)
			certifications.DELETE("/:id", certificationHandler.Delete)
		}

		// Skill routes
		skills := v1.Group("/skills")
		{
			skills.POST("", skillHandler.Create)
			skills.POST("/bulk", skillHandler.CreateBulk)
			skills.GET("", skillHandler.GetAll)
			skills.GET("/:id", skillHandler.GetByID)
			skills.PUT("/:id", skillHandler.Update)
			skills.DELETE("/:id", skillHandler.Delete)
		}

		// Career Preference routes
		preferences := v1.Group("/career-preference")
		{
			preferences.POST("", preferenceHandler.CreateOrUpdateCareerPreference)
			preferences.GET("", preferenceHandler.GetCareerPreference)
		}

		// Position Preference routes
		positionPrefs := v1.Group("/position-preferences")
		{
			positionPrefs.POST("", preferenceHandler.CreatePositionPreferences)
			positionPrefs.GET("", preferenceHandler.GetPositionPreferences)
			positionPrefs.DELETE("/:id", preferenceHandler.DeletePositionPreference)
		}

		// CV routes
		cv := v1.Group("/cv")
		{
			cv.POST("", cvHandler.Upload)
			cv.GET("", cvHandler.Get)
			cv.DELETE("", cvHandler.Delete)
		}

		// Badge routes
		badges := v1.Group("/badges")
		{
			badges.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Badge endpoint - coming soon",
				})
			})
		}
	}

	// Start server
	log.Printf("🚀 %s started on port %s\n", cfg.ServiceName, cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
