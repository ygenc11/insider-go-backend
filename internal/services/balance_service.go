package services

import (
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"time"
)

// Bakiye ekleme/güncelleme
func AddOrUpdateBalance(userID int, amount float64) error {
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil || balance == nil {
		// Bakiye yoksa oluştur
		newBalance := &models.Balance{
			UserID:      userID,
			Amount:      amount,
			LastUpdated: time.Now(),
		}
		return database.CreateBalance(newBalance)
	}

	// Var olan bakiyeyi güncelle
	balance.Amount += amount
	return database.UpdateBalance(userID, balance.Amount)
}

// Kullanıcı bakiyesi çekme
func GetUserBalance(userID int) (float64, error) {
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil || balance == nil {
		return 0, err
	}
	return balance.Amount, nil
}
