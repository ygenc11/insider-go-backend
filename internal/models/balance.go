package models

import (
	"encoding/json"
	"time"
)

type Balance struct {
	UserID      int       `db:"user_id" json:"user_id"`
	Amount      float64   `db:"amount" json:"amount"`
	LastUpdated time.Time `db:"last_updated_at" json:"last_updated_at"`
}

// JSON helper’ları
func (b *Balance) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Balance) FromJSON(data []byte) error {
	return json.Unmarshal(data, b)
}
