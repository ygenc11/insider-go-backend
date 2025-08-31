package services

import (
	"errors"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"time"
)

// Para transferi
func Transfer(fromUserID, toUserID int, amount float64) error {
	// 1. Gönderen bakiyesi
	fromBalance, err := database.GetBalanceByUserID(fromUserID)
	if err != nil {
		return errors.New("gönderen bakiyesi bulunamadı")
	}

	if fromBalance.Amount < amount {
		return errors.New("yetersiz bakiye")
	}

	// Alıcı bakiyesi
	toBalance, err := database.GetBalanceByUserID(toUserID)
	if err != nil {
		return errors.New("alıcı bakiyesi bulunamadı")
	}

	// Bakiyeleri güncelle
	fromBalance.Amount -= amount
	toBalance.Amount += amount

	err = database.UpdateBalance(fromUserID, fromBalance.Amount)
	if err != nil {
		return err
	}

	err = database.UpdateBalance(toUserID, toBalance.Amount)
	if err != nil {
		// rollback (basit)
		database.UpdateBalance(fromUserID, fromBalance.Amount+amount)
		return err
	}

	// Transaction kaydı
	tx := &models.Transaction{
		FromUser:  fromUserID,
		ToUser:    toUserID,
		Amount:    amount,
		Type:      "transfer",
		Status:    "completed",
		CreatedAt: time.Now(),
	}

	return database.CreateTransaction(tx)
}
