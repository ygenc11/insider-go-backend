package services

import (
	"errors"
	"fmt"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"time"
)

// Credit: kullanıcı bakiyesine para ekler ve transaction kaydı oluşturur
func Credit(userID int, amount float64) (float64, error) {
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		return 0, errors.New("balance not found")
	}
	balance.Amount += amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		return 0, err
	}
	tx := &models.Transaction{FromUser: userID, ToUser: userID, Amount: amount, Type: "credit", Status: "completed", CreatedAt: time.Now()}
	if err := database.CreateTransaction(tx); err != nil {
		return 0, err
	}

	// audit log
	_ = LogAction("transaction", tx.ID, "credit", "Credited amount: "+fmt.Sprintf("%.2f", amount))

	return balance.Amount, nil
}

// Debit: kullanıcı bakiyesinden para düşer ve transaction kaydı oluşturur
func Debit(userID int, amount float64) (float64, error) {
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		return 0, errors.New("balance not found")
	}
	if balance.Amount < amount {
		return balance.Amount, errors.New("insufficient funds")
	}
	balance.Amount -= amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		return 0, err
	}
	tx := &models.Transaction{FromUser: userID, ToUser: userID, Amount: amount, Type: "debit", Status: "completed", CreatedAt: time.Now()}
	if err := database.CreateTransaction(tx); err != nil {
		return 0, err
	}

	// audit log
	_ = LogAction("transaction", tx.ID, "debit", "Debited amount: "+fmt.Sprintf("%.2f", amount))

	return balance.Amount, nil
}

// Para transferi: iki bakiye arasında aktarım yapar, transaction kaydı oluşturur; yeni bakiyeleri döner
func Transfer(fromUserID, toUserID int, amount float64) (fromNew float64, toNew float64, err error) {

	fromBalance, err := database.GetBalanceByUserID(fromUserID)
	if err != nil {
		return 0, 0, errors.New("sender balance not found")
	}
	if fromBalance.Amount < amount {
		return fromBalance.Amount, 0, errors.New("insufficient funds")
	}
	toBalance, err := database.GetBalanceByUserID(toUserID)
	if err != nil {
		return 0, 0, errors.New("recipient balance not found")
	}

	fromBalance.Amount -= amount
	toBalance.Amount += amount

	if err := database.UpdateBalance(fromUserID, fromBalance.Amount); err != nil {
		return 0, 0, err
	}
	if err := database.UpdateBalance(toUserID, toBalance.Amount); err != nil {
		// rollback (basit)
		_ = database.UpdateBalance(fromUserID, fromBalance.Amount+amount)
		return 0, 0, err
	}

	tx := &models.Transaction{FromUser: fromUserID, ToUser: toUserID, Amount: amount, Type: "transfer", Status: "completed", CreatedAt: time.Now()}
	if err := database.CreateTransaction(tx); err != nil {
		return 0, 0, err
	}

	// audit log
	_ = LogAction("transaction", tx.ID, "transfer", fmt.Sprintf("Transferred amount: %.2f from user %d to user %d", amount, fromUserID, toUserID))

	return fromBalance.Amount, toBalance.Amount, nil
}

// Sorgular
func GetTransactionsByUser(userID int) ([]*models.Transaction, error) {
	return database.GetTransactionsByUser(userID)
}

func GetTransactionByID(id int) (*models.Transaction, error) {
	return database.GetTransactionByID(id)
}
