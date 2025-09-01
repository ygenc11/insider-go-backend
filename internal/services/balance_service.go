package services

import (
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"log/slog"
	"time"
)

// Bakiye ekleme/güncelleme
func AddOrUpdateBalance(userID int, amount float64) error {
	slog.Info("service.balance.add_or_update.start", "user_id", userID, "delta", amount)
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil || balance == nil {
		// Bakiye yoksa oluştur
		newBalance := &models.Balance{
			UserID:      userID,
			Amount:      amount,
			LastUpdated: time.Now(),
		}
		if err := database.CreateBalance(newBalance); err != nil {
			slog.Error("service.balance.create_failed", "user_id", userID, "err", err)
			return err
		}
		slog.Info("service.balance.created", "user_id", userID, "amount", amount)
		return nil
	}

	// Var olan bakiyeyi güncelle
	balance.Amount += amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		slog.Error("service.balance.update_failed", "user_id", userID, "err", err)
		return err
	}
	slog.Info("service.balance.updated", "user_id", userID, "amount", balance.Amount)
	return nil
}

// Kullanıcı bakiyesi çekme
func GetUserBalance(userID int) (float64, error) {
	slog.Debug("service.balance.get_user_balance", "user_id", userID)
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil || balance == nil {
		return 0, err
	}
	return balance.Amount, nil
}

// GetBalance: kullanıcı bakiyesini getirir
func GetBalance(userID int) (*models.Balance, error) {
	slog.Debug("service.balance.get", "user_id", userID)
	return database.GetBalanceByUserID(userID)
}

// SetBalance: kullanıcı bakiyesini belirli bir değere ayarlar (varsa günceller, yoksa oluşturur)
func SetBalance(userID int, amount float64) (*models.Balance, error) {
	slog.Info("service.balance.set", "user_id", userID, "amount", amount)
	b, err := database.GetBalanceByUserID(userID)
	if err != nil || b == nil {
		newBalance := &models.Balance{
			UserID:      userID,
			Amount:      amount,
			LastUpdated: time.Now(),
		}
		if err := database.CreateBalance(newBalance); err != nil {
			slog.Error("service.balance.create_failed", "user_id", userID, "err", err)
			return nil, err
		}
		return newBalance, nil
	}
	if err := database.UpdateBalance(userID, amount); err != nil {
		slog.Error("service.balance.update_failed", "user_id", userID, "err", err)
		return nil, err
	}
	// Güncellenmiş bakiyeyi tekrar çek
	return database.GetBalanceByUserID(userID)
}

// CalculateBalanceAt: belirli bir zamandaki bakiyeyi hesaplar
func CalculateBalanceAt(userID int, at time.Time) (float64, error) {
	slog.Info("service.balance.calculate_at.start", "user_id", userID, "at", at)
	txs, err := database.GetTransactionsByUser(userID)
	if err != nil {
		slog.Error("service.balance.calculate_at.fetch_failed", "user_id", userID, "err", err)
		return 0, err
	}
	bal := 0.0
	for _, tx := range txs {
		if tx.CreatedAt.After(at) {
			continue
		}
		switch tx.Type {
		case "credit":
			bal += tx.Amount
		case "debit":
			bal -= tx.Amount
		case "transfer":
			if tx.FromUser == userID {
				bal -= tx.Amount
			}
			if tx.ToUser == userID {
				bal += tx.Amount
			}
		}
	}
	slog.Info("service.balance.calculate_at.success", "user_id", userID, "balance", bal)
	return bal, nil
}
