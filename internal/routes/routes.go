package routes

import (
	"insider-go-backend/internal/handlers"
	"insider-go-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Auth endpoints (token gerektirmez)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.RegisterHandler)
			auth.POST("/login", handlers.LoginHandler)
			auth.POST("/refresh", handlers.RefreshHandler)
		}

		// User endpoints (admin rolü gerekli)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
		{
			users.GET("", handlers.GetUsersHandler)
			users.GET("/:id", handlers.GetUserHandler)
			users.PUT("/:id", handlers.UpdateUserHandler)
			users.DELETE("/:id", handlers.DeleteUserHandler)
		}

		// Transaction endpoints (auth gerekli)
		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware())
		{
			transactions.POST("/credit", handlers.CreditHandler)
			transactions.POST("/debit", handlers.DebitHandler)
			transactions.POST("/transfer", handlers.TransferHandler)
			transactions.GET("/history", handlers.TransactionHistoryHandler)
			transactions.GET("/:id", handlers.GetTransactionHandler)
		}

		// Balance endpoints (auth gerekli)
		balances := api.Group("/balances")
		balances.Use(middleware.AuthMiddleware())
		{
			balances.GET("/current", handlers.CurrentBalanceHandler)
			balances.GET("/historical", handlers.HistoricalBalanceHandler)
			balances.GET("/at-time", handlers.BalanceAtTimeHandler)
			// balances.POST("/create", handlers.CreateBalanceHandler) // opsiyonel
		}
	}
}
