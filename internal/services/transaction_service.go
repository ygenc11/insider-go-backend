package services

import (
	"fmt"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"log/slog"
)

// Credit: kullanıcı bakiyesine para ekler ve transaction kaydı oluşturur
func Credit(userID int, amount float64) (float64, error) {
	slog.Info("service.credit.start", "user_id", userID, "amount", amount)
	newBal, tx, err := database.CreditAtomic(userID, amount)
	if err != nil {
		if err.Error() == "balance not found" {
			slog.Error("service.credit.balance_not_found", "user_id", userID, "err", err)
		} else {
			slog.Error("service.credit.failed", "user_id", userID, "err", err)
		}
		return 0, err
	}
	// audit log
	_ = LogAction("transaction", tx.ID, "credit", "Credited amount: "+fmt.Sprintf("%.2f", amount))
	slog.Info("service.credit.success", "user_id", userID, "new_balance", newBal)
	return newBal, nil
}

// Debit: kullanıcı bakiyesinden para düşer ve transaction kaydı oluşturur
func Debit(userID int, amount float64) (float64, error) {
	slog.Info("service.debit.start", "user_id", userID, "amount", amount)
	newBal, tx, err := database.DebitAtomic(userID, amount)
	if err != nil {
		if err.Error() == "insufficient funds" {
			slog.Warn("service.debit.insufficient_funds", "user_id", userID, "amount", amount)
		} else if err.Error() == "balance not found" {
			slog.Error("service.debit.balance_not_found", "user_id", userID, "err", err)
		} else {
			slog.Error("service.debit.failed", "user_id", userID, "err", err)
		}
		return 0, err
	}
	_ = LogAction("transaction", tx.ID, "debit", "Debited amount: "+fmt.Sprintf("%.2f", amount))
	slog.Info("service.debit.success", "user_id", userID, "new_balance", newBal)
	return newBal, nil
}

// Para transferi: iki bakiye arasında aktarım yapar, transaction kaydı oluşturur; yeni bakiyeleri döner
func Transfer(fromUserID, toUserID int, amount float64) (fromNew float64, toNew float64, err error) {
	slog.Info("service.transfer.start", "from_user_id", fromUserID, "to_user_id", toUserID, "amount", amount)
	fromNew, toNew, tx, err := database.TransferAtomic(fromUserID, toUserID, amount)
	if err != nil {
		switch err.Error() {
		case "insufficient funds":
			slog.Warn("service.transfer.insufficient_funds", "from_user_id", fromUserID, "amount", amount)
		case "sender balance not found":
			slog.Error("service.transfer.sender_balance_not_found", "from_user_id", fromUserID, "err", err)
		case "recipient balance not found":
			slog.Error("service.transfer.recipient_balance_not_found", "to_user_id", toUserID, "err", err)
		default:
			slog.Error("service.transfer.failed", "from_user_id", fromUserID, "to_user_id", toUserID, "err", err)
		}
		return 0, 0, err
	}
	_ = LogAction("transaction", tx.ID, "transfer", fmt.Sprintf("Transferred amount: %.2f from user %d to user %d", amount, fromUserID, toUserID))
	slog.Info("service.transfer.success", "from_user_id", fromUserID, "to_user_id", toUserID, "from_new", fromNew, "to_new", toNew)
	return fromNew, toNew, nil
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
