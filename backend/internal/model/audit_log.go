package model

import "time"

// AuditLog 操作审计日志。
type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Action       string    `gorm:"size:64;not null" json:"action"`
	ResourceType string    `gorm:"size:32" json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	Detail       string    `gorm:"size:500" json:"detail"`
	IP           string    `gorm:"size:64" json:"ip"`
	RequestID    string    `gorm:"size:64" json:"request_id"`
	CreatedAt    time.Time `json:"created_at"`
}
