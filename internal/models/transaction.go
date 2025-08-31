package models

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID        int       `gorm:"column:id;primaryKey" db:"id" json:"id"`
	FromUser  int       `gorm:"column:from_user_id" db:"from_user_id" json:"from_user_id"`
	ToUser    int       `gorm:"column:to_user_id" db:"to_user_id" json:"to_user_id"`
	Amount    float64   `gorm:"column:amount" db:"amount" json:"amount"`
	Type      string    `gorm:"column:type" db:"type" json:"type"`       // örn: "credit", "debit", "transfer"
	Status    string    `gorm:"column:status" db:"status" json:"status"` // örn: "pending", "completed", "failed"
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" db:"created_at" json:"created_at"`
}

// JSON helper’ları
func (t *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}
