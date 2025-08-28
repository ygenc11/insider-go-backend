package main

import (
	"fmt"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// DB bağlantısı
	database.ConnectDB("./data.db")

	// Gin router
	r := gin.Default()

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

	fmt.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
