package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"insider-go-backend/internal/database"
	"insider-go-backend/internal/logging"
	mw "insider-go-backend/internal/middleware"
	"insider-go-backend/internal/processor"
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

	// Gin router (manual middlewares: no default Recovery/Logger; we use our own)
	r := gin.New()
	// Init logging and add request/perf/recovery middlewares
	logging.Init()
	r.Use(mw.RequestID())
	r.Use(mw.Recovery())
	r.Use(mw.PerformanceMonitor())
	r.Use(mw.SecurityHeaders())
	r.Use(mw.RequestLogger())
	// IP bazlı rate limit (env ile ayarlanabilir)
	rps := getenvFloat("RATE_LIMIT_RPS", 10)
	burst := getenvFloat("RATE_LIMIT_BURST", 20)
	r.Use(mw.RateLimiter(mw.RateLimiterConfig{RefillRatePerSec: rps, Burst: burst}))

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

	// Transaction processor (opsiyonel): env ile aç/kapa
	if getenv("TXPROC_ENABLED", "true") == "true" {
		workers := getenvInt("TXPROC_WORKERS", 4)
		qcap := getenvInt("TXPROC_QUEUE", 256)
		processor.StartDefault(workers, qcap)
		log.Printf("Transaction processor started (workers=%d, queue=%d)", workers, qcap)
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
	// işlemciyi durdur
	processor.StopDefault()
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

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvFloat(k string, def float64) float64 {
	if v := os.Getenv(k); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
