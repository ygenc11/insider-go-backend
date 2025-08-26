package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreditHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "credit transaction"})
}

func DebitHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "debit transaction"})
}

func TransferHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "transfer transaction"})
}

func TransactionHistoryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "transaction history"})
}

func GetTransactionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get transaction by id"})
}
