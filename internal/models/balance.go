package models

import (
	"encoding/json"
	"time"
)

type Balance struct {
	UserID      int       `gorm:"column:user_id;primaryKey;autoIncrement:false;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" db:"user_id" json:"user_id"`
	Amount      float64   `gorm:"column:amount;type:numeric(18,2);default:0" db:"amount" json:"amount"`
	LastUpdated time.Time `gorm:"column:last_updated_at;autoUpdateTime" db:"last_updated_at" json:"last_updated_at"`
}

// JSON helper’ları
func (b *Balance) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Balance) FromJSON(data []byte) error {
	return json.Unmarshal(data, b)
}
