package database

import (
	"insider-go-backend/internal/models"
	"time"
)

// Yeni transaction ekle
func CreateTransaction(tx *models.Transaction) error {
	tx.CreatedAt = time.Now()
	query := `INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, tx.FromUser, tx.ToUser, tx.Amount, tx.Type, tx.Status, tx.CreatedAt)
	return err
}

// Kullanıcıya ait tüm transactionları getir
func GetTransactionsByUser(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := `SELECT * FROM transactions WHERE from_user_id = ? OR to_user_id = ?`
	err := DB.Select(&transactions, query, userID, userID)
	return transactions, err
}

// ID ile transaction getir
func GetTransactionByID(id int) (*models.Transaction, error) {
	var tx models.Transaction
	query := `SELECT * FROM transactions WHERE id = ?`
	err := DB.Get(&tx, query, id)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
