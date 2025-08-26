package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

var DB *sqlx.DB

// ConnectDB: Veritabanı bağlantısını kurar ve DB değişkenini initialize eder
func ConnectDB(dsn string) {
	var err error
	DB, err = sqlx.Connect("sqlite3", dsn)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	fmt.Println("✅ Database connected")
}
