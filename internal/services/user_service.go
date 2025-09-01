package services

import (
	"errors"
	"log/slog"
	"time"

	"insider-go-backend/internal/config"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT config from env (fallback defaults inside config)
var jwtCfg = config.GetJWT()

// Kullanıcı kaydı
func RegisterUser(username, email, password, role string) (*models.User, error) {
	slog.Info("service.user.register.start", "username", username, "email", email, "role", role)
	existingUser, _ := database.GetUserByEmail(email)
	if existingUser != nil {
		slog.Warn("service.user.register.email_exists", "email", email)
		return nil, errors.New("email already registered")
	}

	existingUser, _ = database.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := database.CreateUser(user); err != nil {
		slog.Error("service.user.register.create_failed", "email", email, "err", err)
		return nil, err
	}

	// audit log for user creation
	_ = LogAction("user", user.ID, "create", "new user registered")

	slog.Info("service.user.register.success", "user_id", user.ID)
	return user, nil
}

// Kullanıcı giriş + JWT üretimi
func AuthenticateUser(email, password string) (string, error) {
	slog.Info("service.user.login.start", "email", email)
	user, err := database.GetUserByEmail(email)
	if err != nil || user == nil {
		slog.Warn("service.user.login.user_not_found", "email", email)
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// JWT oluştur
	tokenString, err := GenerateJWT(user)
	if err != nil {
		slog.Error("service.user.login.jwt_failed", "user_id", user.ID, "err", err)
		return "", err
	}

	// audit log for user login
	_ = LogAction("user", user.ID, "login", "user logged in")

	slog.Info("service.user.login.success", "user_id", user.ID)
	return tokenString, nil
}

// JWT üretme
func GenerateJWT(user *models.User) (string, error) {
	slog.Debug("service.user.jwt.generate", "user_id", user.ID)
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"user_email": user.Email,
		"role":       user.Role,
		"exp":        time.Now().Add(jwtCfg.AccessTTL).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		slog.Error("service.user.jwt.sign_failed", "user_id", user.ID, "err", err)
		return "", err
	}
	return s, nil
}

// JWT çözme
func ParseJWT(tokenStr string) (userID int, role string, err error) {
	slog.Debug("service.user.jwt.parse.start")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// HS256 doğrulama
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtCfg.Secret), nil
	})

	if err != nil {
		slog.Warn("service.user.jwt.parse.error", "err", err)
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			slog.Warn("service.user.jwt.claims.invalid_user_id")
			return 0, "", errors.New("invalid user_id in token")
		}
		userID = int(userIDFloat)

		role, ok = claims["role"].(string)
		if !ok {
			slog.Warn("service.user.jwt.claims.invalid_role")
			return 0, "", errors.New("invalid role in token")
		}

		return userID, role, nil
	}

	slog.Warn("service.user.jwt.claims.invalid")
	return 0, "", errors.New("invalid token claims")
}

// Rol kontrolü
func CheckUserRole(user *models.User, role string) bool {
	return user.Role == role
}

// create balance for user
func CreateBalanceForUser(userID int, initialAmount float64) error {
	slog.Info("service.user.create_balance", "user_id", userID, "initial_amount", initialAmount)
	balance := &models.Balance{
		UserID:      userID,
		Amount:      initialAmount,
		LastUpdated: time.Now(),
	}

	// audit log for balance creation
	_ = LogAction("balance", userID, "create", "initial balance created for user")

	if err := database.CreateBalance(balance); err != nil {
		slog.Error("service.user.create_balance_failed", "user_id", userID, "err", err)
		return err
	}
	return nil
}

// Users domain services (handlers bu katmanı kullanmalı)
func ListUsers() ([]*models.User, error) {
	slog.Info("service.user.list")
	_ = LogAction("user", 0, "list", "list all users")
	return database.GetAllUsers()
}

func GetUser(id int) (*models.User, error) {
	slog.Info("service.user.get", "user_id", id)
	_ = LogAction("user", id, "get", "get user details")
	return database.GetUserByID(id)
}

func UpdateUser(id int, username, email, role string) error {
	slog.Info("service.user.update", "user_id", id)
	_ = LogAction("user", id, "update", "update user details")
	return database.UpdateUser(id, username, email, role)
}

func DeleteUser(id int) error {
	slog.Info("service.user.delete", "user_id", id)
	_ = LogAction("user", id, "delete", "delete user")
	return database.DeleteUser(id)
}
