package main

import (
	"fmt"
	"os"
	"time"

	"insider-go-backend/internal/database"
	"insider-go-backend/internal/logging"
	mw "insider-go-backend/internal/middleware"
	"insider-go-backend/internal/routes"

	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	// DB connection
	dsn := getenv("DB_DSN", "./data.db")
	database.ConnectDB(dsn)

	// Gin router
	r := gin.Default()
	// Init logging and add request logging middleware
	logging.Init()
	r.Use(mw.RequestLogger())

	// Middleware: CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Tüm route’ları kaydet
	routes.RegisterRoutes(r)

	port := getenv("PORT", "8080")
	addr := ":" + port

	// timeoutlarla birlikte HTTP server
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server başlat
	go func() {
		fmt.Printf("Server running at http://localhost:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, shutting down server...")

	shutdownTimeout := getdur("SHUTDOWN_TIMEOUT", 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	// log dosyasını kapat
	logging.Close()
	log.Println("Server gracefully stopped")
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getdur(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
