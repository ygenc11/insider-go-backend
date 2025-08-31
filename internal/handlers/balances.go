package handlers

import (
	"net/http"
	"time"

	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Kullanıcının mevcut bakiyesi
func CurrentBalanceHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "balance not found"})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// Tarihsel bakiye
func HistoricalBalanceHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	transactions, err := database.GetTransactionsByUser(userID)
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

	atTime, err := time.Parse(time.RFC3339, timeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	transactions, err := database.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	balance := 0.0
	for _, tx := range transactions {
		if tx.CreatedAt.After(atTime) {
			continue
		}
		switch tx.Type {
		case "credit":
			balance += tx.Amount
		case "debit":
			balance -= tx.Amount
		case "transfer":
			if tx.FromUser == userID {
				balance -= tx.Amount
			}
			if tx.ToUser == userID {
				balance += tx.Amount
			}
		}
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

	balance := &models.Balance{
		UserID: userID,
		Amount: req.Amount,
	}

	if err := database.CreateBalance(balance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "balance created", "balance": balance})
}
