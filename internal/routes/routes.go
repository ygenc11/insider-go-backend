package routes

import (
	"insider-go-backend/internal/handlers"
	"insider-go-backend/internal/middleware"
	"insider-go-backend/internal/processor"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(r *gin.Engine) {
	// Metrics endpoint for Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
		}

		// Ops: işlemci kuyruğu ve istatistik (admin rolü gerekli olabilir)
		ops := api.Group("/ops")
		ops.Use(middleware.AuthMiddleware())
		{
			ops.GET("/perf", func(c *gin.Context) {
				stats := middleware.GetPerfStats()
				c.JSON(200, stats)
			})

			ops.POST("/enqueue", func(c *gin.Context) {
				var r struct {
					Op       string  `json:"op"`
					UserID   int     `json:"user_id"`
					ToUserID int     `json:"to_user_id"`
					Amount   float64 `json:"amount"`
				}
				if err := c.ShouldBindJSON(&r); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				p := processor.GetDefault()
				if p == nil {
					c.JSON(503, gin.H{"error": "processor not running"})
					return
				}
				job := processor.TxJob{Op: processor.TxOp(r.Op), UserID: r.UserID, ToUserID: r.ToUserID, Amount: r.Amount}
				if ok := p.TryEnqueue(job); !ok {
					c.JSON(429, gin.H{"error": "queue full"})
					return
				}
				c.JSON(202, gin.H{"status": "enqueued"})
			})

			ops.GET("/stats", func(c *gin.Context) {
				p := processor.GetDefault()
				if p == nil {
					c.JSON(503, gin.H{"error": "processor not running"})
					return
				}
				enq, proc, ok, fail := p.Stats()
				c.JSON(200, gin.H{"enqueued": enq, "processed": proc, "succeeded": ok, "failed": fail})
			})
		}
	}
}
