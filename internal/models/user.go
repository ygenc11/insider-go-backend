package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `gorm:"column:id;primaryKey" db:"id" json:"id"`
	Username  string    `gorm:"column:username" db:"username" json:"username"`
	Email     string    `gorm:"column:email" db:"email" json:"email"`
	Password  string    `gorm:"column:password_hash" db:"password_hash" json:"-"`
	Role      string    `gorm:"column:role" db:"role" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" db:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" db:"updated_at" json:"updated_at"`
}

func (u *User) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) FromJSON(data []byte) error {
	return json.Unmarshal(data, u)
}
