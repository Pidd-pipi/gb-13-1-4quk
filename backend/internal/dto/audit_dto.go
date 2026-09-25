package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// AuditQuery 审计日志查询参数。
type AuditQuery struct {
	UserID       uint   `form:"user_id"`
	Action       string `form:"action"`
	ResourceType string `form:"resource_type"`
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
}

// AuditLogDTO 审计日志响应体。
type AuditLogDTO struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   uint   `json:"resource_id"`
	Detail       string `json:"detail"`
	IP           string `json:"ip"`
	RequestID    string `json:"request_id"`
	CreatedAt    string `json:"created_at"`
}

// FromAuditLog converts a model.AuditLog to AuditLogDTO.
func FromAuditLog(a *model.AuditLog) AuditLogDTO {
	return AuditLogDTO{
		ID:           a.ID,
		UserID:       a.UserID,
		Action:       a.Action,
		ResourceType: a.ResourceType,
		ResourceID:   a.ResourceID,
		Detail:       a.Detail,
		IP:           a.IP,
		RequestID:    a.RequestID,
		CreatedAt:    util.FormatTime(a.CreatedAt),
	}
}
