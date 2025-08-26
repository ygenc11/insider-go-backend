package database

import (
	"insider-go-backend/internal/models"
	"time"
)

// Yeni log ekle
func InsertAuditLog(log *models.AuditLog) error {
	log.CreatedAt = time.Now()
	query := `INSERT INTO audit_logs (entity_type, entity_id, action, details, created_at)
	          VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, log.Entity, log.EntityID, log.Action, log.Details, log.CreatedAt)
	return err
}

// Logları getir
func GetAuditLogsByEntity(entity string, entityID int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := `SELECT * FROM audit_logs WHERE entity_type = ? AND entity_id = ?`
	err := DB.Select(&logs, query, entity, entityID)
	return logs, err
}
