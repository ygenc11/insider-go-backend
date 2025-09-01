package database

import (
	"insider-go-backend/internal/models"

	"gorm.io/gorm"
)

type gormAuditLogRepository struct{ db *gorm.DB }

func NewGormAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &gormAuditLogRepository{db: db}
}

func (r *gormAuditLogRepository) InsertAuditLog(log *models.AuditLog) error {
	return r.db.Table("audit_logs").Create(log).Error
}

func (r *gormAuditLogRepository) GetAuditLogsByEntity(entity string, entityID int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Table("audit_logs").Where("entity_type = ? AND entity_id = ?", entity, entityID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *gormAuditLogRepository) GetAllAuditLogs() ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Table("audit_logs").Order("created_at DESC").Find(&logs).Error
	return logs, err
}
