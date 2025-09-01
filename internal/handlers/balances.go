package handlers

import (
	"net/http"
	"time"

	"insider-go-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// Kullanıcının mevcut bakiyesi
func CurrentBalanceHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	balance, err := services.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "balance not found"})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// Tarihsel bakiye
func HistoricalBalanceHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	transactions, err := services.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	var historical []map[string]interface{}
	for _, tx := range transactions {
		historical = append(historical, map[string]interface{}{
			"id":         tx.ID,
			"from_user":  tx.FromUser,
			"to_user":    tx.ToUser,
			"amount":     tx.Amount,
			"type":       tx.Type,
			"status":     tx.Status,
			"created_at": tx.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, historical)
}

// Belirli bir zamanda bakiye
func BalanceAtTimeHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	timeParam := c.Query("at_time")

	// Accept RFC3339 or date-only (YYYY-MM-DD). Date-only is treated as end of that day (UTC).
	atTime, err := time.Parse(time.RFC3339, timeParam)
	if err != nil {
		if d, derr := time.Parse("2006-01-02", timeParam); derr == nil {
			// end of day UTC: 23:59:59
			atTime = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, time.UTC)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format; use RFC3339 or YYYY-MM-DD"})
			return
		}
	}

	balance, err := services.CalculateBalanceAt(userID, atTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute balance at time"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance_at_time": balance, "at_time": atTime})
}

// Yeni bakiye oluştur
func CreateBalanceHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := services.SetBalance(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "balance created", "balance": updated})
}
