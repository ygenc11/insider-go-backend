package database

import (
	"insider-go-backend/internal/models"
	"sync"

	"gorm.io/gorm"
)

var balanceMutex = &sync.RWMutex{}

// Kullanıcı bakiyesi getir
func GetBalanceByUserID(userID int) (*models.Balance, error) {
	var balance models.Balance
	if err := DB.Table("balances").Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return nil, err
	}
	return &balance, nil
}

// Bakiyeyi thread-safe güncelle
func UpdateBalance(userID int, amount float64) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	return DB.Table("balances").Where("user_id = ?", userID).Updates(map[string]interface{}{
		"amount":          amount,
		"last_updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}).Error
}

// Yeni bakiye ekle
func CreateBalance(balance *models.Balance) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	return DB.Table("balances").Create(balance).Error
}

// Bakiyeyi artır veya azalt (thread-safe)
func AdjustBalance(userID int, delta float64) error {
	balanceMutex.Lock()
	defer balanceMutex.Unlock()
	var balance models.Balance
	if err := DB.Table("balances").Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return err
	}
	newAmount := balance.Amount + delta
	return DB.Table("balances").Where("user_id = ?", userID).Updates(map[string]interface{}{
		"amount":          newAmount,
		"last_updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}).Error
}
