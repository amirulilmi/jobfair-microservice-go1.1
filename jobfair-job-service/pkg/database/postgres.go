package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	
	// Retry connection up to 10 times with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// Get underlying SQL DB
			sqlDB, err := db.DB()
			if err == nil {
				// Connection pool settings
				sqlDB.SetMaxIdleConns(10)
				sqlDB.SetMaxOpenConns(100)
				sqlDB.SetConnMaxLifetime(time.Hour)

				// Test connection
				if err := sqlDB.Ping(); err == nil {
					log.Println("✅ Database connected successfully")
					return db, nil
				}
			}
		}
		
		waitTime := time.Duration(i+1) * 2 * time.Second
		log.Printf("⏳ Failed to connect to database (attempt %d/%d). Retrying in %v... Error: %v", 
			i+1, maxRetries, waitTime, err)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
