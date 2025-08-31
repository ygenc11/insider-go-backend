package database

import (
	"insider-go-backend/internal/models"
)

// Yeni transaction ekle
func CreateTransaction(tx *models.Transaction) error {
	return DB.Table("transactions").Create(tx).Error
}

// Kullanıcıya ait tüm transactionları getir
func GetTransactionsByUser(userID int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := DB.Table("transactions").
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Find(&transactions).Error
	return transactions, err
}

// ID ile transaction getir
func GetTransactionByID(id int) (*models.Transaction, error) {
	var tx models.Transaction
	if err := DB.Table("transactions").First(&tx, id).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}
