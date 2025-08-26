package services

import (
	"insider-go-backend/internal/database"
	"insider-go-backend/internal/models"
	"time"
)

// Yeni log ekle
func LogAction(entity string, entityID int, action, details string) error {
	logEntry := &models.AuditLog{
		Entity:    entity,
		EntityID:  entityID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return database.InsertAuditLog(logEntry)
}

// Logları getir
func GetEntityLogs(entity string, entityID int) ([]models.AuditLog, error) {
	return database.GetAuditLogsByEntity(entity, entityID)
}
