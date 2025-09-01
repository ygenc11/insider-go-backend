package services

import (
	"insider-go-backend/internal/models"
	"time"
)

// UserService arayüzü
type UserService interface {
	RegisterUser(username, email, password, role string) (*models.User, error)
	AuthenticateUser(email, password string) (access string, refresh string, err error)
	GenerateJWT(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
	RefreshAccessToken(refreshToken string) (string, error)
	ParseJWT(tokenStr string) (userID int, role string, err error)
	CheckUserRole(user *models.User, role string) bool
	CreateBalanceForUser(userID int, initialAmount float64) error
	ListUsers() ([]*models.User, error)
	GetUser(id int) (*models.User, error)
	UpdateUser(id int, username, email, role string) error
	DeleteUser(id int) error
}

// BalanceService arayüzü
type BalanceService interface {
	AddOrUpdateBalance(userID int, amount float64) error
	GetUserBalance(userID int) (float64, error)
	GetBalance(userID int) (*models.Balance, error)
	SetBalance(userID int, amount float64) (*models.Balance, error)
	CalculateBalanceAt(userID int, at time.Time) (float64, error)
}

// TransactionService arayüzü
type TransactionService interface {
	Credit(userID int, amount float64) (float64, error)
	Debit(userID int, amount float64) (float64, error)
	Transfer(fromUserID, toUserID int, amount float64) (fromNew float64, toNew float64, err error)
	GetTransactionsByUser(userID int) ([]*models.Transaction, error)
	GetTransactionByID(id int) (*models.Transaction, error)
}

// AuditLogService arayüzü
type AuditLogService interface {
	LogAction(entity string, entityID int, action, details string) error
	GetEntityLogs(entity string, entityID int) ([]models.AuditLog, error)
}
