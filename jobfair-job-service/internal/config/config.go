package config

import (
	"os"
)

type Config struct {
	DatabaseURL       string
	Port              string
	JWTSecret         string
	RabbitMQURL       string
	AuthServiceURL    string
	CompanyServiceURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://jobfair_user:jobfair_pass@localhost:5435/jobfair_jobs?sslmode=disable"),
		Port:              getEnv("PORT", "8082"),
		JWTSecret:         getEnv("JWT_SECRET", "default-secret"),
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
		CompanyServiceURL: getEnv("COMPANY_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
