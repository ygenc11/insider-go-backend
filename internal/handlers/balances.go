package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
)

// Kullanıcının mevcut bakiyesi
func CurrentBalanceHandler(c *gin.Context) {
	userIDParam := c.Query("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "balance not found"})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// Kullanıcının tarihsel bakiyeleri (tüm transactionlar üzerinden)
func HistoricalBalanceHandler(c *gin.Context) {
	userIDParam := c.Query("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

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

// Kullanıcının belirli bir zamanda bakiyesi
func BalanceAtTimeHandler(c *gin.Context) {
	userIDParam := c.Query("user_id")
	timeParam := c.Query("at_time") // ISO8601 format: 2025-08-26T17:00:00Z

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "balance_at_time": balance, "at_time": atTime})
}

// Yeni bakiye oluştur
func CreateBalanceHandler(c *gin.Context) {
	var req struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance := &models.Balance{
		UserID: req.UserID,
		Amount: req.Amount,
	}

	if err := database.CreateBalance(balance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "balance created", "balance": balance})
}
