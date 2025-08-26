package main

import (
	"fmt"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/handlers"
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

	// Routes
	api := r.Group("/api/v1")
	{
		// Auth endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.RegisterHandler)
			auth.POST("/login", handlers.LoginHandler)
			auth.POST("/refresh", handlers.RefreshHandler)
		}

		// User endpoints
		users := api.Group("/users")
		{
			users.GET("", handlers.GetUsersHandler)
			users.GET("/:id", handlers.GetUserHandler)
			users.PUT("/:id", handlers.UpdateUserHandler)
			users.DELETE("/:id", handlers.DeleteUserHandler)
		}

		// Transaction endpoints
		transactions := api.Group("/transactions")
		{
			transactions.POST("/credit", handlers.CreditHandler)
			transactions.POST("/debit", handlers.DebitHandler)
			transactions.POST("/transfer", handlers.TransferHandler)
			transactions.GET("/history", handlers.TransactionHistoryHandler)
			transactions.GET("/:id", handlers.GetTransactionHandler)
		}

		// Balance endpoints
		balances := api.Group("/balances")
		{
			balances.GET("/current", handlers.CurrentBalanceHandler)
			balances.GET("/historical", handlers.HistoricalBalanceHandler)
			balances.GET("/at-time", handlers.BalanceAtTimeHandler)
		}
	}

	fmt.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
