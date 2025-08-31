package database

import (
	"insider-go-backend/internal/models"

	"gorm.io/gorm"
)

// Tüm kullanıcıları getir
func GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := DB.Table("users").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ID’ye göre kullanıcı getir
func GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := DB.Table("users").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Email’e göre kullanıcı bul
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := DB.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Username’a göre kullanıcı bul
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := DB.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Yeni kullanıcı ekle
func CreateUser(user *models.User) error {
	return DB.Table("users").Create(user).Error
}

// Kullanıcıyı güncelle
func UpdateUser(id int, username, email, role string) error {
	updates := map[string]interface{}{
		"username":   username,
		"email":      email,
		"role":       role,
		"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}
	return DB.Table("users").Where("id = ?", id).Updates(updates).Error
}

// Kullanıcıyı sil
func DeleteUser(id int) error {
	return DB.Table("users").Where("id = ?", id).Delete(&models.User{}).Error
}
