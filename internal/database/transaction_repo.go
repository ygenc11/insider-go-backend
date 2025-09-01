package database

import (
	"errors"
	"insider-go-backend/internal/models"
	"time"

	"gorm.io/gorm"
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

// CreditAtomic: belirli kullanıcının bakiyesini atomik olarak artırır ve aynı transaction içinde transaction kaydı oluşturur
func CreditAtomic(userID int, amount float64) (float64, *models.Transaction, error) {
	var newAmount float64
	rec := &models.Transaction{}
	err := DB.Transaction(func(tx *gorm.DB) error {
		// bakiye var mı?
		var b models.Balance
		if err := tx.Table("balances").Where("user_id = ?", userID).First(&b).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("balance not found")
			}
			return err
		}
		// artır
		if err := tx.Exec("UPDATE balances SET amount = amount + ?, last_updated_at = CURRENT_TIMESTAMP WHERE user_id = ?", amount, userID).Error; err != nil {
			return err
		}
		if err := tx.Table("balances").Select("amount").Where("user_id = ?", userID).Scan(&newAmount).Error; err != nil {
			return err
		}
		*rec = models.Transaction{FromUser: userID, ToUser: userID, Amount: amount, Type: "credit", Status: "completed", CreatedAt: time.Now()}
		if err := tx.Table("transactions").Create(rec).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return newAmount, rec, nil
}

// DebitAtomic: bakiyeden atomik olarak düşer (yeterli bakiye şartı) ve transaction kaydı oluşturur
func DebitAtomic(userID int, amount float64) (float64, *models.Transaction, error) {
	var newAmount float64
	rec := &models.Transaction{}
	err := DB.Transaction(func(tx *gorm.DB) error {
		// bakiye var mı?
		var b models.Balance
		if err := tx.Table("balances").Where("user_id = ?", userID).First(&b).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("balance not found")
			}
			return err
		}
		// yeterli bakiye ile azalt
		res := tx.Exec("UPDATE balances SET amount = amount - ?, last_updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND amount >= ?", amount, userID, amount)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("insufficient funds")
		}
		if err := tx.Table("balances").Select("amount").Where("user_id = ?", userID).Scan(&newAmount).Error; err != nil {
			return err
		}
		*rec = models.Transaction{FromUser: userID, ToUser: userID, Amount: amount, Type: "debit", Status: "completed", CreatedAt: time.Now()}
		if err := tx.Table("transactions").Create(rec).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return newAmount, rec, nil
}

// TransferAtomic: iki kullanıcı arasında atomik transfer yapar ve transaction kaydı oluşturur
func TransferAtomic(fromUserID, toUserID int, amount float64) (float64, float64, *models.Transaction, error) {
	var fromAmt, toAmt float64
	rec := &models.Transaction{}
	err := DB.Transaction(func(tx *gorm.DB) error {
		// bakiye var mı? (her iki taraf)
		var fromB, toB models.Balance
		if err := tx.Table("balances").Where("user_id = ?", fromUserID).First(&fromB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("sender balance not found")
			}
			return err
		}
		if err := tx.Table("balances").Where("user_id = ?", toUserID).First(&toB).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("recipient balance not found")
			}
			return err
		}
		// gönderenden düş (yeterli bakiye ile)
		res := tx.Exec("UPDATE balances SET amount = amount - ?, last_updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND amount >= ?", amount, fromUserID, amount)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("insufficient funds")
		}
		// alıcıya ekle
		if err := tx.Exec("UPDATE balances SET amount = amount + ?, last_updated_at = CURRENT_TIMESTAMP WHERE user_id = ?", amount, toUserID).Error; err != nil {
			return err
		}
		if err := tx.Table("balances").Select("amount").Where("user_id = ?", fromUserID).Scan(&fromAmt).Error; err != nil {
			return err
		}
		if err := tx.Table("balances").Select("amount").Where("user_id = ?", toUserID).Scan(&toAmt).Error; err != nil {
			return err
		}
		*rec = models.Transaction{FromUser: fromUserID, ToUser: toUserID, Amount: amount, Type: "transfer", Status: "completed", CreatedAt: time.Now()}
		if err := tx.Table("transactions").Create(rec).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, 0, nil, err
	}
	return fromAmt, toAmt, rec, nil
}
