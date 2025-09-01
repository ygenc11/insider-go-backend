package services

import (
	"insider-go-backend/internal/models"
	"time"
)

// Varsayılan servis örnekleri (fonksiyonları mevcut global fonksiyonlara delege eder)
var (
	defaultUserService        UserService        = userServiceImpl{}
	defaultBalanceService     BalanceService     = balanceServiceImpl{}
	defaultTransactionService TransactionService = transactionServiceImpl{}
	defaultAuditLogService    AuditLogService    = auditLogServiceImpl{}
)

// Getter'lar
func UserSvc() UserService               { return defaultUserService }
func BalanceSvc() BalanceService         { return defaultBalanceService }
func TransactionSvc() TransactionService { return defaultTransactionService }
func AuditLogSvc() AuditLogService       { return defaultAuditLogService }

// Setters (test veya özel implementasyonlar için)
func SetUserSvc(s UserService)               { defaultUserService = s }
func SetBalanceSvc(s BalanceService)         { defaultBalanceService = s }
func SetTransactionSvc(s TransactionService) { defaultTransactionService = s }
func SetAuditLogSvc(s AuditLogService)       { defaultAuditLogService = s }

// Basit implementasyonlar: varolan paket-level fonksiyonlara delege
type userServiceImpl struct{}

func (userServiceImpl) RegisterUser(username, email, password, role string) (*models.User, error) {
	return RegisterUser(username, email, password, role)
}
func (userServiceImpl) AuthenticateUser(email, password string) (string, string, error) {
	return AuthenticateUser(email, password)
}
func (userServiceImpl) GenerateJWT(user *models.User) (string, error) { return GenerateJWT(user) }
func (userServiceImpl) GenerateRefreshToken(user *models.User) (string, error) {
	return GenerateRefreshToken(user)
}
func (userServiceImpl) RefreshAccessToken(refreshToken string) (string, error) {
	return RefreshAccessToken(refreshToken)
}
func (userServiceImpl) ParseJWT(tokenStr string) (int, string, error) { return ParseJWT(tokenStr) }
func (userServiceImpl) CheckUserRole(user *models.User, role string) bool {
	return CheckUserRole(user, role)
}
func (userServiceImpl) CreateBalanceForUser(userID int, initialAmount float64) error {
	return CreateBalanceForUser(userID, initialAmount)
}
func (userServiceImpl) ListUsers() ([]*models.User, error)   { return ListUsers() }
func (userServiceImpl) GetUser(id int) (*models.User, error) { return GetUser(id) }
func (userServiceImpl) UpdateUser(id int, username, email, role string) error {
	return UpdateUser(id, username, email, role)
}
func (userServiceImpl) DeleteUser(id int) error { return DeleteUser(id) }

type balanceServiceImpl struct{}

func (balanceServiceImpl) AddOrUpdateBalance(userID int, amount float64) error {
	return AddOrUpdateBalance(userID, amount)
}
func (balanceServiceImpl) GetUserBalance(userID int) (float64, error)     { return GetUserBalance(userID) }
func (balanceServiceImpl) GetBalance(userID int) (*models.Balance, error) { return GetBalance(userID) }
func (balanceServiceImpl) SetBalance(userID int, amount float64) (*models.Balance, error) {
	return SetBalance(userID, amount)
}
func (balanceServiceImpl) CalculateBalanceAt(userID int, at time.Time) (float64, error) {
	return CalculateBalanceAt(userID, at)
}

type transactionServiceImpl struct{}

func (transactionServiceImpl) Credit(userID int, amount float64) (float64, error) {
	return Credit(userID, amount)
}
func (transactionServiceImpl) Debit(userID int, amount float64) (float64, error) {
	return Debit(userID, amount)
}
func (transactionServiceImpl) Transfer(fromUserID, toUserID int, amount float64) (float64, float64, error) {
	return Transfer(fromUserID, toUserID, amount)
}
func (transactionServiceImpl) GetTransactionsByUser(userID int) ([]*models.Transaction, error) {
	return GetTransactionsByUser(userID)
}
func (transactionServiceImpl) GetTransactionByID(id int) (*models.Transaction, error) {
	return GetTransactionByID(id)
}

type auditLogServiceImpl struct{}

func (auditLogServiceImpl) LogAction(entity string, entityID int, action, details string) error {
	return LogAction(entity, entityID, action, details)
}
func (auditLogServiceImpl) GetEntityLogs(entity string, entityID int) ([]models.AuditLog, error) {
	return GetEntityLogs(entity, entityID)
}
