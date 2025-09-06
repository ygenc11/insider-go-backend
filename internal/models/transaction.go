package models

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID        int       `gorm:"column:id;primaryKey" db:"id" json:"id"`
	FromUser  int       `gorm:"column:from_user_id;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" db:"from_user_id" json:"from_user_id"`
	ToUser    int       `gorm:"column:to_user_id;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" db:"to_user_id" json:"to_user_id"`
	Amount    float64   `gorm:"column:amount;type:numeric(18,2)" db:"amount" json:"amount"`
	Type      string    `gorm:"column:type;index" db:"type" json:"type"`
	Status    string    `gorm:"column:status;index" db:"status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index" db:"created_at" json:"created_at"`
}

// JSON helper’ları
func (t *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) FromJSON(data []byte) error {
	return json.Unmarshal(data, t)
}
