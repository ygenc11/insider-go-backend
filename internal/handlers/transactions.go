package handlers

import (
	"net/http"
	"strconv"

	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Credit işlemi
func CreditHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "balance not found"})
		return
	}

	balance.Amount += req.Amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update balance"})
		return
	}

	tx := &models.Transaction{
		FromUser: userID,
		ToUser:   userID,
		Amount:   req.Amount,
		Type:     "credit",
		Status:   "completed",
	}
	if err := database.CreateTransaction(tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "credit successful", "balance": balance.Amount})
}

// Debit işlemi
func DebitHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "balance not found"})
		return
	}

	if balance.Amount < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	balance.Amount -= req.Amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update balance"})
		return
	}

	tx := &models.Transaction{
		FromUser: userID,
		ToUser:   userID,
		Amount:   req.Amount,
		Type:     "debit",
		Status:   "completed",
	}
	if err := database.CreateTransaction(tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "debit successful", "balance": balance.Amount})
}

// Transfer işlemi
func TransferHandler(c *gin.Context) {
	fromUser := c.GetInt("user_id")

	var req struct {
		ToUser int     `json:"to_user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromBalance, err := database.GetBalanceByUserID(fromUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sender balance not found"})
		return
	}
	if fromBalance.Amount < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	toBalance, err := database.GetBalanceByUserID(req.ToUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "receiver balance not found"})
		return
	}

	fromBalance.Amount -= req.Amount
	toBalance.Amount += req.Amount

	if err := database.UpdateBalance(fromUser, fromBalance.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sender balance"})
		return
	}
	if err := database.UpdateBalance(req.ToUser, toBalance.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update receiver balance"})
		return
	}

	tx := &models.Transaction{
		FromUser: fromUser,
		ToUser:   req.ToUser,
		Amount:   req.Amount,
		Type:     "transfer",
		Status:   "completed",
	}
	if err := database.CreateTransaction(tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "transfer successful",
		"from_balance": fromBalance.Amount,
		"to_balance":   toBalance.Amount,
	})
}

// Transaction geçmişi
func TransactionHistoryHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	transactions, err := database.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// ID’ye göre transaction getir
func GetTransactionHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	tx, err := database.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}
