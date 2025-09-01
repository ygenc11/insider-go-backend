package services

import (
	"errors"
	"fmt"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"log/slog"
	"time"
)

// Credit: kullanıcı bakiyesine para ekler ve transaction kaydı oluşturur
func Credit(userID int, amount float64) (float64, error) {
	slog.Info("service.credit.start", "user_id", userID, "amount", amount)
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		slog.Error("service.credit.balance_not_found", "user_id", userID, "err", err)
		return 0, errors.New("balance not found")
	}
	balance.Amount += amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		slog.Error("service.credit.update_balance_failed", "user_id", userID, "err", err)
		return 0, err
	}
	tx := &models.Transaction{FromUser: userID, ToUser: userID, Amount: amount, Type: "credit", Status: "completed", CreatedAt: time.Now()}
	if err := database.CreateTransaction(tx); err != nil {
		slog.Error("service.credit.create_tx_failed", "user_id", userID, "err", err)
		return 0, err
	}

	// audit log
	_ = LogAction("transaction", tx.ID, "credit", "Credited amount: "+fmt.Sprintf("%.2f", amount))

	slog.Info("service.credit.success", "user_id", userID, "new_balance", balance.Amount)
	return balance.Amount, nil
}

// Debit: kullanıcı bakiyesinden para düşer ve transaction kaydı oluşturur
func Debit(userID int, amount float64) (float64, error) {
	slog.Info("service.debit.start", "user_id", userID, "amount", amount)
	balance, err := database.GetBalanceByUserID(userID)
	if err != nil {
		slog.Error("service.debit.balance_not_found", "user_id", userID, "err", err)
		return 0, errors.New("balance not found")
	}
	if balance.Amount < amount {
		slog.Warn("service.debit.insufficient_funds", "user_id", userID, "balance", balance.Amount, "amount", amount)
		return balance.Amount, errors.New("insufficient funds")
	}
	balance.Amount -= amount
	if err := database.UpdateBalance(userID, balance.Amount); err != nil {
		slog.Error("service.debit.update_balance_failed", "user_id", userID, "err", err)
		return 0, err
	}
	tx := &models.Transaction{FromUser: userID, ToUser: userID, Amount: amount, Type: "debit", Status: "completed", CreatedAt: time.Now()}
	if err := database.CreateTransaction(tx); err != nil {
		slog.Error("service.debit.create_tx_failed", "user_id", userID, "err", err)
		return 0, err
	}

	// audit log
	_ = LogAction("transaction", tx.ID, "debit", "Debited amount: "+fmt.Sprintf("%.2f", amount))

	slog.Info("service.debit.success", "user_id", userID, "new_balance", balance.Amount)
	return balance.Amount, nil
}

// Para transferi: iki bakiye arasında aktarım yapar, transaction kaydı oluşturur; yeni bakiyeleri döner
func Transfer(fromUserID, toUserID int, amount float64) (fromNew float64, toNew float64, err error) {
	slog.Info("service.transfer.start", "from_user_id", fromUserID, "to_user_id", toUserID, "amount", amount)

	fromBalance, err := database.GetBalanceByUserID(fromUserID)
	if err != nil {
		slog.Error("service.transfer.sender_balance_not_found", "from_user_id", fromUserID, "err", err)
		return 0, 0, errors.New("sender balance not found")
	}
	if fromBalance.Amount < amount {
		slog.Warn("service.transfer.insufficient_funds", "from_user_id", fromUserID, "balance", fromBalance.Amount, "amount", amount)
		return fromBalance.Amount, 0, errors.New("insufficient funds")
	}
	toBalance, err := database.GetBalanceByUserID(toUserID)
	if err != nil {
		slog.Error("service.transfer.recipient_balance_not_found", "to_user_id", toUserID, "err", err)
		return 0, 0, errors.New("recipient balance not found")
	}

	fromBalance.Amount -= amount
	toBalance.Amount += amount

	if err := database.UpdateBalance(fromUserID, fromBalance.Amount); err != nil {
		slog.Error("service.transfer.update_from_failed", "from_user_id", fromUserID, "err", err)
		return 0, 0, err
	}
	if err := database.UpdateBalance(toUserID, toBalance.Amount); err != nil {
		// rollback (basit)
		_ = database.UpdateBalance(fromUserID, fromBalance.Amount+amount)
		slog.Error("service.transfer.update_to_failed", "to_user_id", toUserID, "err", err)
		return 0, 0, err
	}

	tx := &models.Transaction{FromUser: fromUserID, ToUser: toUserID, Amount: amount, Type: "transfer", Status: "completed", CreatedAt: time.Now()}
	if err := database.CreateTransaction(tx); err != nil {
		slog.Error("service.transfer.create_tx_failed", "from_user_id", fromUserID, "to_user_id", toUserID, "err", err)
		return 0, 0, err
	}

	// audit log
	_ = LogAction("transaction", tx.ID, "transfer", fmt.Sprintf("Transferred amount: %.2f from user %d to user %d", amount, fromUserID, toUserID))

	slog.Info("service.transfer.success", "from_user_id", fromUserID, "to_user_id", toUserID, "from_new", fromBalance.Amount, "to_new", toBalance.Amount)
	return fromBalance.Amount, toBalance.Amount, nil
}

// Sorgular
func GetTransactionsByUser(userID int) ([]*models.Transaction, error) {
	slog.Info("service.transactions.list_by_user", "user_id", userID)
	return database.GetTransactionsByUser(userID)
}

func GetTransactionByID(id int) (*models.Transaction, error) {
	slog.Info("service.transactions.get_by_id", "id", id)
	return database.GetTransactionByID(id)
}
