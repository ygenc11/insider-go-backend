package database

import (
	"fmt"
	"log"

	"insider-go-backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection using GORM
func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	fmt.Println("✅ Database connected (GORM)")

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Balance{},
		&models.AuditLog{},
	); err != nil {
		log.Printf("AutoMigrate failed: %v", err)
	} else {
		log.Printf("AutoMigrate applied (DEV)")
	}
}
