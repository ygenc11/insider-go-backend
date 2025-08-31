package database

import (
	"insider-go-backend/internal/models"
)

// Yeni log ekle
func InsertAuditLog(log *models.AuditLog) error {
	return DB.Table("audit_logs").Create(log).Error
}

// Belirli entity ve ID’ye ait logları getir
func GetAuditLogsByEntity(entity string, entityID int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := DB.Table("audit_logs").Where("entity_type = ? AND entity_id = ?", entity, entityID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// Tüm logları getir
func GetAllAuditLogs() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := DB.Table("audit_logs").Order("created_at DESC").Find(&logs).Error
	return logs, err
}
