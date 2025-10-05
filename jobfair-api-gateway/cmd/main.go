package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServiceConfig struct {
	Name string
	URL  string
	Path string
}

func main() {
	// Load environment variables
	port := getEnv("PORT", "8000")
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8080")
	companyServiceURL := getEnv("COMPANY_SERVICE_URL", "http://localhost:8081")
	jobServiceURL := getEnv("JOB_SERVICE_URL", "http://localhost:8082")
	userProfileServiceURL := getEnv("USER_PROFILE_SERVICE_URL", "http://localhost:8083")

	// Initialize Gin router
	router := gin.Default()

	// Disable automatic trailing slash redirect to prevent 301 loops
	router.RedirectTrailingSlash = false

	// Set max multipart memory for file uploads (50MB)
	router.MaxMultipartMemory = 50 << 20

	// CORS middleware
	router.Use(corsMiddleware())

	// Logging middleware
	router.Use(loggingMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "api-gateway",
			"timestamp": time.Now().Format(time.RFC3339),
			"services": gin.H{
				"auth":    authServiceURL,
				"company": companyServiceURL,
				"job":     jobServiceURL,
				"profile": userProfileServiceURL,
			},
		})
	})

	// Service status endpoint
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"gateway": "running",
			"version": "1.0.0",
			"uptime":  time.Now().Format(time.RFC3339),
		})
	})

	// Create reverse proxies
	authProxy := createReverseProxy(authServiceURL, "auth-service")
	companyProxy := createReverseProxy(companyServiceURL, "company-service")
	jobProxy := createReverseProxy(jobServiceURL, "job-service")
	profileProxy := createReverseProxy(userProfileServiceURL, "user-profile-service")

	// ==================== AUTH SERVICE ROUTES ====================
	// Login & Token routes
	router.Any("/api/v1/login", proxyHandler(authProxy, "/api/v1/login"))
	router.Any("/api/v1/refresh", proxyHandler(authProxy, "/api/v1/refresh"))

	// Get current user (protected)
	router.Any("/api/v1/me", proxyHandler(authProxy, "/api/v1/me"))

	// Registration routes (multi-step)
	router.Any("/api/v1/register/*proxyPath", proxyHandler(authProxy, "/api/v1/register"))
	router.Any("/api/v1/register", proxyHandler(authProxy, "/api/v1/register"))

	// Other auth routes
	router.Any("/api/v1/auth/*proxyPath", proxyHandler(authProxy, "/api/v1/auth"))

	// ==================== STATIC FILES - AUTH SERVICE ====================
	// Profile photos and company logos
	router.Any("/uploads/profiles/*filepath", proxyHandler(authProxy, "/uploads/profiles"))
	router.Any("/uploads/companies/*filepath", proxyHandler(authProxy, "/uploads/companies"))

	// ==================== COMPANY SERVICE ROUTES ====================
	// Company CRUD routes
	router.Any("/api/v1/companies/*proxyPath", proxyHandler(companyProxy, "/api/v1/companies"))
	router.Any("/api/v1/companies", proxyHandler(companyProxy, "/api/v1/companies"))

	// My company route
	router.Any("/api/v1/my-company", proxyHandler(companyProxy, "/api/v1/my-company"))

	// ==================== JOB SERVICE ROUTES ====================
	router.Any("/api/v1/jobs/*proxyPath", proxyHandler(jobProxy, "/api/v1/jobs"))
	router.Any("/api/v1/jobs", proxyHandler(jobProxy, "/api/v1/jobs"))
	router.Any("/api/v1/applications/*proxyPath", proxyHandler(jobProxy, "/api/v1/applications"))
	router.Any("/api/v1/applications", proxyHandler(jobProxy, "/api/v1/applications"))

	// ==================== USER PROFILE SERVICE ROUTES ====================
	// Profile routes
	router.Any("/api/v1/profiles/*proxyPath", proxyHandler(profileProxy, "/api/v1/profiles"))
	router.Any("/api/v1/profiles", proxyHandler(profileProxy, "/api/v1/profiles"))

	// Work Experience routes
	router.Any("/api/v1/work-experiences/*proxyPath", proxyHandler(profileProxy, "/api/v1/work-experiences"))
	router.Any("/api/v1/work-experiences", proxyHandler(profileProxy, "/api/v1/work-experiences"))

	// Education routes
	router.Any("/api/v1/educations/*proxyPath", proxyHandler(profileProxy, "/api/v1/educations"))
	router.Any("/api/v1/educations", proxyHandler(profileProxy, "/api/v1/educations"))

	// Certifications routes
	router.Any("/api/v1/certifications/*proxyPath", proxyHandler(profileProxy, "/api/v1/certifications"))
	router.Any("/api/v1/certifications", proxyHandler(profileProxy, "/api/v1/certifications"))

	// Skills routes
	router.Any("/api/v1/skills/*proxyPath", proxyHandler(profileProxy, "/api/v1/skills"))
	router.Any("/api/v1/skills", proxyHandler(profileProxy, "/api/v1/skills"))

	// Career Preference routes
	router.Any("/api/v1/career-preference/*proxyPath", proxyHandler(profileProxy, "/api/v1/career-preference"))
	router.Any("/api/v1/career-preference", proxyHandler(profileProxy, "/api/v1/career-preference"))

	// Position Preferences routes
	router.Any("/api/v1/position-preferences/*proxyPath", proxyHandler(profileProxy, "/api/v1/position-preferences"))
	router.Any("/api/v1/position-preferences", proxyHandler(profileProxy, "/api/v1/position-preferences"))

	// CV routes
	router.Any("/api/v1/cv/*proxyPath", proxyHandler(profileProxy, "/api/v1/cv"))
	router.Any("/api/v1/cv", proxyHandler(profileProxy, "/api/v1/cv"))

	// Banner routes
	router.Any("/api/v1/banner/*proxyPath", proxyHandler(profileProxy, "/api/v1/banner"))
	router.Any("/api/v1/banner", proxyHandler(profileProxy, "/api/v1/banner"))

	// Badges routes
	router.Any("/api/v1/badges/*proxyPath", proxyHandler(profileProxy, "/api/v1/badges"))
	router.Any("/api/v1/badges", proxyHandler(profileProxy, "/api/v1/badges"))

	// ==================== STATIC FILES - USER PROFILE SERVICE ====================
	// CV files
	router.Any("/uploads/cv/*filepath", proxyHandler(profileProxy, "/uploads/cv"))
	// Banner images
	router.Any("/uploads/banners/*filepath", proxyHandler(profileProxy, "/uploads/banners"))

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Route not found",
			"path":    c.Request.URL.Path,
		})
	})

	// Start server
	log.Printf("🚀 API Gateway starting on port %s", port)
	log.Printf("📡 Proxying to services:")
	log.Printf("   - Auth Service: %s", authServiceURL)
	log.Printf("   - Company Service: %s", companyServiceURL)
	log.Printf("   - Job Service: %s", jobServiceURL)
	log.Printf("   - User Profile Service: %s", userProfileServiceURL)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start API Gateway: %v", err)
	}
}

// proxyHandler creates a Gin handler that proxies requests
func proxyHandler(proxy *httputil.ReverseProxy, pathPrefix string) gin.HandlerFunc {
	return gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		// Store original path
		originalPath := r.URL.Path

		// Keep the original path without modification
		// The target service will handle the full path including prefix

		// Log the proxy
		log.Printf("🔄 Proxying: %s", originalPath)

		// Serve the request with original path
		proxy.ServeHTTP(w, r)
	})
}

// createReverseProxy creates a reverse proxy for a service
func createReverseProxy(targetURL string, serviceName string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse URL for %s: %v", serviceName, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Custom director to modify the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Origin-Host", target.Host)
		req.Host = target.Host
	}

	// Error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("❌ Proxy error for %s: %v", serviceName, err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"success":false,"message":"Service temporarily unavailable","error":"` + err.Error() + `"}`))
	}

	return proxy
}

// corsMiddleware sets up CORS
func corsMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
		AllowMethods:     strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS,PATCH"), ","),
		AllowHeaders:     strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,Accept,Origin,User-Agent,Cache-Control,Keep-Alive,X-Requested-With"), ","),
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowFiles:       true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}

// loggingMiddleware logs each request
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s - %d (%v)",
			method,
			path,
			c.ClientIP(),
			statusCode,
			latency,
		)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
