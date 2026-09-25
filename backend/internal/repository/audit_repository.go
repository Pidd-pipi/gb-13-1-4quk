package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
)

// AuditRepository handles audit log persistence.
type AuditRepository struct{ db *gorm.DB }

// NewAuditRepository creates an AuditRepository.
func NewAuditRepository(db *gorm.DB) *AuditRepository { return &AuditRepository{db: db} }

// Create inserts an audit log.
func (r *AuditRepository) Create(a *model.AuditLog) error { return translate(r.db.Create(a).Error) }

// List returns audit logs filtered by AuditQuery.
func (r *AuditRepository) List(q dto.AuditQuery) ([]model.AuditLog, int64, error) {
	var items []model.AuditLog
	var total int64
	db := r.db.Model(&model.AuditLog{})
	if q.UserID > 0 {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	}
	if q.ResourceType != "" {
		db = db.Where("resource_type = ?", q.ResourceType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(q.PageSize * (q.Page - 1)).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
