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
	userProfileServiceURL := getEnv("USER_PROFILE_SERVICE_URL", "http://localhost:8083")

	// Initialize Gin router
	router := gin.Default()

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

	// Create reverse proxy for Auth Service
	authProxy := createReverseProxy(authServiceURL, "auth-service")
	router.Any("/api/v1/auth/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		// Remove /api/v1/auth prefix before proxying
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/auth")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		// Add /api/v1/auth back for the target service
		r.URL.Path = "/api/v1/auth" + r.URL.Path
		authProxy.ServeHTTP(w, r)
	}))

	// Create reverse proxy for Company Service
	companyProxy := createReverseProxy(companyServiceURL, "company-service")
	router.Any("/api/v1/companies/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/companies")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/companies" + r.URL.Path
		companyProxy.ServeHTTP(w, r)
	}))

	// Create reverse proxy for User Profile Service
	profileProxy := createReverseProxy(userProfileServiceURL, "user-profile-service")
	router.Any("/api/v1/profiles/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/profiles")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/profiles" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Work Experience routes
	router.Any("/api/v1/work-experiences/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/work-experiences")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/work-experiences" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Education routes
	router.Any("/api/v1/educations/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/educations")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/educations" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Certifications routes
	router.Any("/api/v1/certifications/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/certifications")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/certifications" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Skills routes
	router.Any("/api/v1/skills/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/skills")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/skills" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Career Preference routes
	router.Any("/api/v1/career-preference/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/career-preference")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/career-preference" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Position Preferences routes
	router.Any("/api/v1/position-preferences/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/position-preferences")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/position-preferences" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// CV routes
	router.Any("/api/v1/cv/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/cv")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/cv" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

	// Badges routes
	router.Any("/api/v1/badges/*proxyPath", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1/badges")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		r.URL.Path = "/api/v1/badges" + r.URL.Path
		profileProxy.ServeHTTP(w, r)
	}))

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
	log.Printf("   - User Profile Service: %s", userProfileServiceURL)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start API Gateway: %v", err)
	}
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
		AllowMethods:     strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
		AllowHeaders:     strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"), ","),
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
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
