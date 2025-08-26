package database

import (
	"insider-go-backend/internal/models"
	"time"
)

// Kullanıcı bakiyesi getir
func GetBalanceByUserID(userID int) (*models.Balance, error) {
	var balance models.Balance
	query := `SELECT * FROM balances WHERE user_id = ?`
	err := DB.Get(&balance, query, userID)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// Bakiyeyi güncelle
func UpdateBalance(userID int, amount float64) error {
	query := `UPDATE balances SET amount = ?, last_updated_at = ? WHERE user_id = ?`
	_, err := DB.Exec(query, amount, time.Now(), userID)
	return err
}

// Yeni bakiye ekle
func CreateBalance(balance *models.Balance) error {
	balance.LastUpdated = time.Now()
	query := `INSERT INTO balances (user_id, amount, last_updated_at) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, balance.UserID, balance.Amount, balance.LastUpdated)
	return err
}
