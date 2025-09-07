package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"insider-go-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Veritabanına bağlan ve otomatik migrate işlemi yap
func ConnectDB(dsn string) {
	// If explicit DSN not provided build from env pieces
	if dsn == "" {
		dsn = buildPostgresDSN()
	}
	if os.Getenv("DEBUG_DB") == "1" {
		log.Printf("Connecting DB with DSN: %s", redactPassword(dsn))
	}
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(getenvInt("DB_MAX_OPEN", 20))
		sqlDB.SetMaxIdleConns(getenvInt("DB_MAX_IDLE", 10))
		sqlDB.SetConnMaxLifetime(getenvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	}
	fmt.Println("Postgres connected (GORM)")

	// Varsayılan repository implementasyonlarını başlat
	InitDefaultRepos(DB)

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

func buildPostgresDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "postgres"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "app"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "app"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "appdb"
	}
	ssl := os.Getenv("DB_SSLMODE")
	if ssl == "" {
		ssl = "disable"
	}
	opts := []string{fmt.Sprintf("host=%s", host), fmt.Sprintf("port=%s", port), fmt.Sprintf("user=%s", user), fmt.Sprintf("password=%s", pass), fmt.Sprintf("dbname=%s", name), fmt.Sprintf("sslmode=%s", ssl)}
	if ext := os.Getenv("DB_EXTRA"); ext != "" {
		opts = append(opts, ext)
	}
	return strings.Join(opts, " ")
}

// redactPassword: mask password segment inside DSN for safe logging
func redactPassword(dsn string) string {
	parts := strings.Split(dsn, " ")
	for i, p := range parts {
		if strings.HasPrefix(p, "password=") {
			parts[i] = "password=***"
		}
	}
	return strings.Join(parts, " ")
}

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
func getenvDuration(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
