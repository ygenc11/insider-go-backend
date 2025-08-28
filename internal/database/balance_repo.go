package database

import (
	"insider-go-backend/internal/models"
	"sync"
	"time"
)

var balanceMutex = &sync.RWMutex{}

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

// Bakiyeyi thread-safe güncelle
func UpdateBalance(userID int, amount float64) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()

	query := `UPDATE balances SET amount = ?, last_updated_at = ? WHERE user_id = ?`
	_, err := DB.Exec(query, amount, time.Now(), userID)
	return err
}

// Yeni bakiye ekle
func CreateBalance(balance *models.Balance) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()

	balance.LastUpdated = time.Now()
	query := `INSERT INTO balances (user_id, amount, last_updated_at) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, balance.UserID, balance.Amount, balance.LastUpdated)
	return err
}

// Bakiyeyi artır veya azalt (thread-safe)
func AdjustBalance(userID int, delta float64) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()

	balance, err := GetBalanceByUserID(userID)
	if err != nil {
		return err
	}

	newAmount := balance.Amount + delta
	return UpdateBalance(userID, newAmount)
}
