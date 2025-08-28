package routes

import (
	"insider-go-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
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
			balances.POST("/create", handlers.CreateBalanceHandler)
		}
	}
}
