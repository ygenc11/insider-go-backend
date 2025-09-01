package database

import (
	"insider-go-backend/internal/models"

	"gorm.io/gorm"
)

type gormUserRepository struct{ db *gorm.DB }

func NewGormUserRepository(db *gorm.DB) UserRepository { return &gormUserRepository{db: db} }

func (r *gormUserRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Table("users").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *gormUserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.Table("users").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) CreateUser(user *models.User) error {
	return r.db.Table("users").Create(user).Error
}

func (r *gormUserRepository) UpdateUser(id int, username, email, role string) error {
	updates := map[string]interface{}{
		"username":   username,
		"email":      email,
		"role":       role,
		"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}
	return r.db.Table("users").Where("id = ?", id).Updates(updates).Error
}

func (r *gormUserRepository) DeleteUser(id int) error {
	return r.db.Table("users").Where("id = ?", id).Delete(&models.User{}).Error
}
