package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = buildDSN()
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// migrations klasörünün yolu
	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations applied successfully (UP).")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations rolled back successfully (DOWN).")
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}
}

func buildDSN() string {
	host := env("DB_HOST", "postgres")
	port := env("DB_PORT", "5432")
	user := env("DB_USER", "app")
	pass := env("DB_PASSWORD", "app")
	name := env("DB_NAME", "appdb")
	ssl := env("DB_SSLMODE", "disable")
	base := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, name, ssl)
	if extra := os.Getenv("DB_EXTRA"); extra != "" {
		base = base + " " + strings.TrimSpace(extra)
	}
	return base
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
