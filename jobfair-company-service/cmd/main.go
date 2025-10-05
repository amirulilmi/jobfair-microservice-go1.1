package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"jobfair-company-service/internal/config"
	"jobfair-company-service/internal/consumers"
	"jobfair-company-service/internal/handlers"
	"jobfair-company-service/internal/middleware"
	"jobfair-company-service/internal/repository"
	"jobfair-company-service/internal/services"
	"jobfair-company-service/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

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

	// 🚀 Initialize Event Consumer
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/" // default
	}

	companyRepo := repository.NewCompanyRepository(db)
	
	eventConsumer, err := consumers.NewCompanyEventConsumer(rabbitmqURL, companyRepo)
	if err != nil {
		log.Fatal("❌ Failed to create event consumer:", err)
	}
	defer eventConsumer.Close()

	// Start consuming events in background
	if err := eventConsumer.Start(); err != nil {
		log.Fatal("❌ Failed to start event consumer:", err)
	}
	log.Println("✅ Event consumer started")

	// Initialize services
	companyService := services.NewCompanyService(companyRepo)

	// Initialize handlers
	companyHandler := handlers.NewCompanyHandler(companyService)

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
			"service": "company-service",
			"version": "1.0.0",
		})
	})

	// Static file serving for uploads
	router.Static("/uploads", "./uploads")

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret"
	}

	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	// API Routes
	api := router.Group("/api/v1")
	{
		// Public routes
		public := api.Group("")
		{
			public.GET("/companies", companyHandler.ListCompanies)
			public.GET("/companies/:id", companyHandler.GetCompany)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(jwtMiddleware)
		{
			// Company management
			protected.GET("/my-company", companyHandler.GetMyCompany)
			protected.POST("/companies", companyHandler.CreateCompany)
			protected.PUT("/companies/:id", companyHandler.UpdateCompany)

			// File uploads
			protected.POST("/companies/:id/logo", companyHandler.UploadLogo)
			protected.POST("/companies/:id/banner", companyHandler.UploadBanner)
			protected.POST("/companies/:id/videos", companyHandler.UploadVideo)
			protected.POST("/companies/:id/gallery", companyHandler.UploadGallery)

			// Analytics
			protected.GET("/companies/:id/analytics", companyHandler.GetAnalytics)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Company service starting on port %s", port)
		log.Printf("📊 Health check: http://localhost:%s/health", port)
		log.Printf("🔑 API endpoint: http://localhost:%s/api/v1", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("🛑 Shutting down server...")
	log.Println("✅ Server exited")
}
