package handlers

import (
	"net/http"
	"strconv"

	"insider-go-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	ToUser int     `json:"to_user_id"`
}

// POST /transactions/credit
func CreditHandler(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt("user_id")
	newBal, err := services.Credit(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "credited", "new_balance": newBal})
}

// POST /transactions/debit
func DebitHandler(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt("user_id")
	newBal, err := services.Debit(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "debited", "new_balance": newBal})
}

// POST /transactions/transfer
func TransferHandler(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fromUserID := c.GetInt("user_id")
	toUserID := req.ToUser

	if toUserID == 0 || toUserID == fromUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipient"})
		return
	}
	fromNew, _, err := services.Transfer(fromUserID, toUserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer completed", "old_balance": fromNew + req.Amount, "new_balance": fromNew, "amount_transferred": req.Amount})
}

// GET /transactions/history
func TransactionHistoryHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	txs, err := services.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": txs})
}

// GET /transactions/:id
func GetTransactionHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	tx, err := services.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}
