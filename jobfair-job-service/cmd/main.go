package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"jobfair-job-service/internal/config"
	"jobfair-job-service/internal/consumers"
	"jobfair-job-service/internal/handlers"
	"jobfair-job-service/internal/middleware"
	"jobfair-job-service/internal/repository"
	"jobfair-job-service/internal/services"
	"jobfair-job-service/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Initialize repositories
	jobRepo := repository.NewJobRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	savedJobRepo := repository.NewSavedJobRepository(db)
	companyRepo := repository.NewCompanyRepository(db)

	// Initialize services
	jobService := services.NewJobService(jobRepo, applicationRepo, savedJobRepo, companyRepo)
	applicationService := services.NewApplicationService(applicationRepo, jobRepo)

	// Initialize handlers
	jobHandler := handlers.NewJobHandler(jobService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	adminHandler := handlers.NewAdminHandler(companyRepo, jobService)

	// Initialize and start event consumer
	companyConsumer, err := consumers.NewCompanyEventConsumer(cfg.RabbitMQURL, companyRepo)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to initialize company event consumer: %v", err)
		log.Println("Service will continue without event consumption")
	} else {
		go func() {
			if err := companyConsumer.Start(); err != nil {
				log.Printf("❌ Company event consumer error: %v", err)
			}
		}()
		log.Println("✅ Company event consumer started")
	}

	// Setup Gin
	router := gin.Default()
	
	// Disable automatic trailing slash redirect to prevent 301 loops
	router.RedirectTrailingSlash = false
	
	router.MaxMultipartMemory = 10 << 20 // 10MB

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
			"service": "job-service",
			"version": "1.0.0",
		})
	})

	// JWT Secret
	jwtSecret := cfg.JWTSecret
	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)
	optionalJWT := middleware.OptionalJWTMiddleware(jwtSecret)

	// API Routes
	api := router.Group("/api/v1")
	{
		// Public routes (with optional JWT for personalization)
		public := api.Group("")
		public.Use(optionalJWT)
		{
			public.GET("/jobs", jobHandler.ListJobs)
			public.GET("/jobs/:id", jobHandler.GetJob)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(jwtMiddleware)
		{
			// Job management (company only)
			protected.POST("/jobs", jobHandler.CreateJob)
			protected.PUT("/jobs/:id", jobHandler.UpdateJob)
			protected.DELETE("/jobs/:id", jobHandler.DeleteJob)
			protected.POST("/jobs/:id/publish", jobHandler.PublishJob)
			protected.POST("/jobs/:id/close", jobHandler.CloseJob)
			protected.GET("/jobs/my", jobHandler.GetMyJobs)

			// Job application (job seeker only)
			protected.POST("/jobs/:id/apply", jobHandler.ApplyToJob)
			protected.POST("/jobs/bulk-apply", jobHandler.BulkApply)

			// Saved jobs (bookmarks)
			protected.POST("/jobs/:id/save", jobHandler.SaveJob)
			protected.DELETE("/jobs/:id/save", jobHandler.UnsaveJob)
			protected.GET("/jobs/saved", jobHandler.GetSavedJobs)

			// Applications management
			protected.GET("/applications/my", applicationHandler.GetMyApplications)
			protected.GET("/applications/:id", applicationHandler.GetApplication)
			protected.PUT("/applications/:id/status", applicationHandler.UpdateApplicationStatus)
			protected.DELETE("/applications/:id", applicationHandler.WithdrawApplication)
			protected.GET("/applications/stats", applicationHandler.GetApplicationStats)

			// Applications by job (company only)
			protected.GET("/jobs/:id/applications", applicationHandler.GetApplicationsByJobID)
		}

		// Admin routes (admin only)
		admin := api.Group("/admin")
		admin.Use(jwtMiddleware)
		{
			// Company mapping management
			admin.POST("/sync-company-mapping", adminHandler.SyncCompanyMapping)
			admin.GET("/company-mappings", adminHandler.GetCompanyMappings)
			
			// Health checks
			admin.GET("/health/data-consistency", adminHandler.HealthCheckDataConsistency)
		}
	}

	port := cfg.Port
	if port == "" {
		port = "8082"
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Job service starting on port %s", port)
		log.Printf("📊 Health check: http://localhost:%s/health", port)
		log.Printf("🔑 API endpoint: http://localhost:%s/api/v1", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("🛑 Shutting down server...")
	
	// Close consumer
	if companyConsumer != nil {
		if err := companyConsumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}
	
	log.Println("✅ Server exited")
}
