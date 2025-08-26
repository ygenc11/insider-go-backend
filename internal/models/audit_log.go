package models

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID        int       `db:"id" json:"id"`
	Entity    string    `db:"entity_type" json:"entity_type"` // örn: "user", "transaction"
	EntityID  int       `db:"entity_id" json:"entity_id"`
	Action    string    `db:"action" json:"action"` // örn: "create", "update", "delete"
	Details   string    `db:"details" json:"details"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// JSON helper’ları
func (a *AuditLog) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AuditLog) FromJSON(data []byte) error {
	return json.Unmarshal(data, a)
}
