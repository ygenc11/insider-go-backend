package database

import (
	"insider-go-backend/internal/models"
	"sync"

	"gorm.io/gorm"
)

type gormBalanceRepository struct {
	db *gorm.DB
	mu *sync.RWMutex
}

func NewGormBalanceRepository(db *gorm.DB) BalanceRepository {
	return &gormBalanceRepository{db: db, mu: &sync.RWMutex{}}
}

func (r *gormBalanceRepository) GetBalanceByUserID(userID int) (*models.Balance, error) {
	var balance models.Balance
	if err := r.db.Table("balances").Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *gormBalanceRepository) UpdateBalance(userID int, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Table("balances").Where("user_id = ?", userID).Updates(map[string]interface{}{
		"amount":          amount,
		"last_updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}).Error
}

func (r *gormBalanceRepository) CreateBalance(balance *models.Balance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Table("balances").Create(balance).Error
}

func (r *gormBalanceRepository) AdjustBalance(userID int, delta float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var balance models.Balance
	if err := r.db.Table("balances").Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return err
	}
	newAmount := balance.Amount + delta
	return r.db.Table("balances").Where("user_id = ?", userID).Updates(map[string]interface{}{
		"amount":          newAmount,
		"last_updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}).Error
}
