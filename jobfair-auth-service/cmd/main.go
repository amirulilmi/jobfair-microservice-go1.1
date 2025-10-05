// File: cmd/main.go
package main

import (
	"log"
	"os"

	"jobfair-auth-service/internal/config"
	"jobfair-auth-service/internal/handlers"
	"jobfair-auth-service/internal/middleware"
	"jobfair-auth-service/internal/repository"
	"jobfair-auth-service/internal/services"
	"jobfair-auth-service/pkg/database"

	"github.com/jobfair/shared/events" // Import shared events library

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 🚀 Initialize Event Publisher
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/" // default
	}

	eventPublisher, err := events.NewPublisher(rabbitmqURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to RabbitMQ:", err)
	}
	defer eventPublisher.Close()
	log.Println("✅ Event publisher initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewJobSeekerProfileRepository(db)
	companyProfileRepo := repository.NewCompanyBasicProfileRepository(db)
	otpRepo := repository.NewOTPRepository(db)

	// Initialize services with event publisher
	registrationService := services.NewRegistrationService(
		userRepo,
		profileRepo,
		companyProfileRepo,
		otpRepo,
		cfg.JWTSecret,
		eventPublisher, // 🎯 Inject event publisher
	)
	authService := services.NewAuthService(userRepo, profileRepo, companyProfileRepo, cfg.JWTSecret)

	// Initialize handlers
	registrationHandler := handlers.NewRegistrationHandler(registrationService)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	router := gin.Default()
	
	// Disable automatic trailing slash redirect to prevent 301 loops
	router.RedirectTrailingSlash = false

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "auth-service",
			"version": "1.0.0",
		})
	})

	// Static file serving for uploads
	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{
		// Registration flow (multi-step) - UNIFIED for both job_seeker and company
		register := api.Group("/register")
		{
			register.POST("/init", registrationHandler.RegisterStep1)                          // Email & Password + UserType
			register.PUT("/profile", authMiddleware, registrationHandler.CompleteBasicProfile) // Unified: detects job_seeker or company
			register.POST("/send-otp", authMiddleware, registrationHandler.SendPhoneOTP)       // Send OTP
			register.POST("/verify-otp", registrationHandler.VerifyPhoneOTP)                   // Verify OTP

			// Job seeker specific steps
			register.POST("/employment", authMiddleware, registrationHandler.SetEmploymentStatus) // Job seeker only
			register.POST("/preferences", authMiddleware, registrationHandler.SetJobPreferences)  // Job seeker only
			// register.POST("/permissions", authMiddleware, registrationHandler.SetPermissions)     // Job seeker only

			// Unified photo/logo upload
			register.POST("/photo", authMiddleware, registrationHandler.UploadProfilePhoto) // Profile photo / Company logo

			// Debug endpoint
			register.GET("/users", authHandler.GetAllUsers)
		}

		// Authentication
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)
		
		// Get current user profile (protected route)
		api.GET("/me", authMiddleware, authHandler.GetCurrentUser)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Auth service starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/health", port)
	log.Printf("🔑 API endpoint: http://localhost:%s/api/v1", port)
	log.Printf("📁 Static files: http://localhost:%s/uploads", port)
	router.Run(":" + port)
}
