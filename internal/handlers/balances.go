package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentBalanceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "current balance"})
}

func HistoricalBalanceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "historical balance"})
}

func BalanceAtTimeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "balance at specific time"})
}
