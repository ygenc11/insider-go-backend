package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"insider-go-backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Veritabanına bağlan ve otomatik migrate işlemi yap
func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	fmt.Println("✅ Database connected (GORM)")

	if shouldAutoMigrate() {
		if err := DB.AutoMigrate(
			&models.User{},
			&models.Transaction{},
			&models.Balance{},
			&models.AuditLog{},
		); err != nil {
			log.Printf("AutoMigrate failed: %v", err)
		} else {
			log.Printf("AutoMigrate applied")
		}
	}
}

func shouldAutoMigrate() bool {
	v := os.Getenv("AUTO_MIGRATE")
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}
