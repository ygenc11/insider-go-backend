package database

import (
	"insider-go-backend/internal/models"

	"gorm.io/gorm"
)

// UserRepository arayüzü
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(id int, username, email, role string) error
	DeleteUser(id int) error
}

// BalanceRepository arayüzü
type BalanceRepository interface {
	GetBalanceByUserID(userID int) (*models.Balance, error)
	UpdateBalance(userID int, amount float64) error
	CreateBalance(balance *models.Balance) error
	AdjustBalance(userID int, delta float64) error
}

// TransactionRepository arayüzü
type TransactionRepository interface {
	CreateTransaction(tx *models.Transaction) error
	GetTransactionsByUser(userID int) ([]*models.Transaction, error)
	GetTransactionByID(id int) (*models.Transaction, error)
	// Atomik para hareketleri
	CreditAtomic(userID int, amount float64) (float64, *models.Transaction, error)
	DebitAtomic(userID int, amount float64) (float64, *models.Transaction, error)
	TransferAtomic(fromUserID, toUserID int, amount float64) (float64, float64, *models.Transaction, error)
}

// AuditLogRepository arayüzü
type AuditLogRepository interface {
	InsertAuditLog(log *models.AuditLog) error
	GetAuditLogsByEntity(entity string, entityID int) ([]models.AuditLog, error)
	GetAllAuditLogs() ([]models.AuditLog, error)
}

// Varsayılan repo örnekleri
var (
	defaultUserRepo        UserRepository
	defaultBalanceRepo     BalanceRepository
	defaultTransactionRepo TransactionRepository
	defaultAuditLogRepo    AuditLogRepository
)

// InitDefaultRepos: uygulama başlangıcında çağrılmalı
func InitDefaultRepos(db *gorm.DB) {
	defaultUserRepo = NewGormUserRepository(db)
	defaultBalanceRepo = NewGormBalanceRepository(db)
	defaultTransactionRepo = NewGormTransactionRepository(db)
	defaultAuditLogRepo = NewGormAuditLogRepository(db)
}

// Getter'lar
func UserRepo() UserRepository               { return defaultUserRepo }
func BalanceRepo() BalanceRepository         { return defaultBalanceRepo }
func TransactionRepo() TransactionRepository { return defaultTransactionRepo }
func AuditLogRepo() AuditLogRepository       { return defaultAuditLogRepo }

// Setters (test veya özel implementasyonlar için)
func SetUserRepo(r UserRepository)               { defaultUserRepo = r }
func SetBalanceRepo(r BalanceRepository)         { defaultBalanceRepo = r }
func SetTransactionRepo(r TransactionRepository) { defaultTransactionRepo = r }
func SetAuditLogRepo(r AuditLogRepository)       { defaultAuditLogRepo = r }
