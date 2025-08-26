package models

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID        int       `db:"id" json:"id"`
	FromUser  int       `db:"from_user_id" json:"from_user_id"`
	ToUser    int       `db:"to_user_id" json:"to_user_id"`
	Amount    float64   `db:"amount" json:"amount"`
	Type      string    `db:"type" json:"type"`     // örn: "credit", "debit", "transfer"
	Status    string    `db:"status" json:"status"` // örn: "pending", "completed", "failed"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// JSON helper’ları
func (t *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}
