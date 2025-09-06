package models

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID        int       `gorm:"column:id;primaryKey" db:"id" json:"id"`
	Entity    string    `gorm:"column:entity_type;index:idx_audit_entity" db:"entity_type" json:"entity_type"`
	EntityID  int       `gorm:"column:entity_id;index:idx_audit_entity" db:"entity_id" json:"entity_id"`
	Action    string    `gorm:"column:action;index" db:"action" json:"action"`
	Details   string    `gorm:"column:details" db:"details" json:"details"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index" db:"created_at" json:"created_at"`
}

// JSON helper’ları
func (a *AuditLog) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AuditLog) FromJSON(data []byte) error {
	return json.Unmarshal(data, a)
}
