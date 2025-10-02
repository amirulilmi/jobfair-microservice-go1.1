package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	Port               string
	ServiceName        string
	MaxFileSize        int64
	AllowedFileTypes   string
	UploadDir          string
	AWSRegion          string
	AWSBucketName      string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AuthServiceURL     string
	CompanyServiceURL  string
}

func Load() *Config {
	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "5242880"), 10, 64)

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		Port:               getEnv("PORT", "8083"),
		ServiceName:        getEnv("SERVICE_NAME", "user-profile-service"),
		MaxFileSize:        maxFileSize,
		AllowedFileTypes:   getEnv("ALLOWED_FILE_TYPES", ".pdf,.doc,.docx,.jpg,.jpeg,.png"),
		UploadDir:          getEnv("UPLOAD_DIR", "./uploads/cv"),
		AWSRegion:          getEnv("AWS_REGION", "ap-southeast-1"),
		AWSBucketName:      getEnv("AWS_BUCKET_NAME", ""),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
		CompanyServiceURL:  getEnv("COMPANY_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
