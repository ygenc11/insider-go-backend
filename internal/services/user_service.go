package services

import (
	"errors"
	"time"

	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key (prod’da env’den alınmalı)
var jwtSecret = []byte("abcdefghijklmnopqrstuvwxyz") // örnek, production’da env kullan

// Token süresi
const tokenExpiry = 24 * time.Hour

// Kullanıcı kaydı
func RegisterUser(username, email, password, role string) (*models.User, error) {
	existingUser, _ := database.GetUserByEmail(email)
	if existingUser != nil {
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
		return nil, err
	}

	// audit log for user creation
	_ = LogAction("user", user.ID, "create", "new user registered")

	return user, nil
}

// Kullanıcı giriş + JWT üretimi
func AuthenticateUser(email, password string) (string, error) {
	user, err := database.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// JWT oluştur
	tokenString, err := GenerateJWT(user)
	if err != nil {
		return "", err
	}

	// audit log for user login
	_ = LogAction("user", user.ID, "login", "user logged in")

	return tokenString, nil
}

// JWT üretme
func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"user_email": user.Email,
		"role":       user.Role,
		"exp":        time.Now().Add(tokenExpiry).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT çözme
func ParseJWT(tokenStr string) (userID int, role string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// HS256 doğrulama
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", errors.New("invalid user_id in token")
		}
		userID = int(userIDFloat)

		role, ok = claims["role"].(string)
		if !ok {
			return 0, "", errors.New("invalid role in token")
		}

		return userID, role, nil
	}

	return 0, "", errors.New("invalid token claims")
}

// Rol kontrolü
func CheckUserRole(user *models.User, role string) bool {
	return user.Role == role
}

// create balance for user
func CreateBalanceForUser(userID int, initialAmount float64) error {
	balance := &models.Balance{
		UserID:      userID,
		Amount:      initialAmount,
		LastUpdated: time.Now(),
	}

	// audit log for balance creation
	_ = LogAction("balance", userID, "create", "initial balance created for user")

	return database.CreateBalance(balance)
}

// Users domain services (handlers bu katmanı kullanmalı)
func ListUsers() ([]*models.User, error) {
	_ = LogAction("user", 0, "list", "list all users")
	return database.GetAllUsers()
}

func GetUser(id int) (*models.User, error) {
	_ = LogAction("user", id, "get", "get user details")
	return database.GetUserByID(id)
}

func UpdateUser(id int, username, email, role string) error {
	_ = LogAction("user", id, "update", "update user details")
	return database.UpdateUser(id, username, email, role)
}

func DeleteUser(id int) error {
	_ = LogAction("user", id, "delete", "delete user")
	return database.DeleteUser(id)
}
