package database

import (
	"insider-go-backend/internal/models"
)

// DB global değişkeni ConnectDB ile initialize edilmiş olmalı.

// Tüm kullanıcıları getir
func GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := DB.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ID’ye göre kullanıcı getir
func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := DB.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Email’e göre kullanıcı bul
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = ?`
	err := DB.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Yeni kullanıcı ekle
func CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
              VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := DB.Exec(query, user.Username, user.Email, user.Password, user.Role)
	return err
}

// Kullanıcıyı güncelle
func UpdateUser(id int, username, email, role string) error {
	query := `UPDATE users SET username = ?, email = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := DB.Exec(query, username, email, role, id)
	return err
}

// Kullanıcıyı sil
func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}
