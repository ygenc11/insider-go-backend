package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
)

// Kullanıcı kaydı
func RegisterUser(username, email, password, role string) error {
	// Email zaten kayıtlı mı kontrol et
	existingUser, _ := database.GetUserByEmail(email)
	if existingUser != nil {
		return errors.New("email zaten kayıtlı")
	}

	// Şifre hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// User struct oluştur
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	// DB’ye kaydet
	return database.CreateUser(user)
}

// Kullanıcı giriş
func AuthenticateUser(email, password string) (*models.User, error) {
	user, err := database.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	// Şifre doğrula
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("şifre yanlış")
	}

	return user, nil
}

// Rol kontrolü
func CheckUserRole(user *models.User, role string) bool {
	return user.Role == role
}
