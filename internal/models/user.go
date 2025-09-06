package models

import (
	"encoding/json"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID        int       `gorm:"column:id;primaryKey" db:"id" json:"id"`
	Username  string    `gorm:"column:username;uniqueIndex" db:"username" json:"username"`
	Email     string    `gorm:"column:email;uniqueIndex" db:"email" json:"email"`
	Password  string    `gorm:"column:password_hash" db:"password_hash" json:"-"`
	Role      string    `gorm:"column:role;default:user" db:"role" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" db:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" db:"updated_at" json:"updated_at"`
}

func (u *User) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) FromJSON(data []byte) error {
	return json.Unmarshal(data, u)
}

// Alan: izin verilen roller
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// Rolün izin verilen rollerden biri olup olmadığını kontrol eder
func IsValidRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

var (
	// Kullanıcı adı için desen: 3-50 karakter, harf/rakam/_ . -
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,50}$`)
)

// Validate, User için temel alan kurallarını doğrular.
func (u *User) Validate() error {
	var problems []string

	// username: zorunlu, 3-50 karakter, sınırlı karakter kümesi
	if s := strings.TrimSpace(u.Username); s == "" || !usernameRe.MatchString(s) {
		problems = append(problems, "invalid username (3-50 chars, letters/digits/_ . -)")
	}

	// email: zorunlu, net/mail ile RFC'ye uygunluk kontrolü
	if e := strings.TrimSpace(u.Email); e == "" {
		problems = append(problems, "email is required")
	} else if _, err := mail.ParseAddress(e); err != nil {
		problems = append(problems, "invalid email format")
	}

	// role: zorunlu, izin verilen kümede olmalı
	if r := strings.TrimSpace(u.Role); r == "" || !IsValidRole(r) {
		problems = append(problems, "invalid role (allowed: user, admin)")
	}

	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}
