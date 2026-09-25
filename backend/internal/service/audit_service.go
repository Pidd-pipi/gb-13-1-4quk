package service

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/util"
)

// AuditService records and queries operation audit logs.
type AuditService struct {
	auditRepo *repository.AuditRepository
	logger    *slog.Logger
}

// NewAuditService creates an AuditService.
func NewAuditService(auditRepo *repository.AuditRepository, logger *slog.Logger) *AuditService {
	return &AuditService{auditRepo: auditRepo, logger: logger}
}

// Record writes an audit log entry.
func (s *AuditService) Record(userID uint, action, resourceType string, resourceID uint, detail, ip, requestID string) {
	entry := &model.AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detail,
		IP:           ip,
		RequestID:    requestID,
	}
	if err := s.auditRepo.Create(entry); err != nil {
		s.logger.Warn(constants.LogAuditRecorded, "error", err)
		return
	}
	s.logger.Info(constants.LogAuditRecorded, "user_id", userID, "action", action, "resource", resourceType+"/"+fmt.Sprint(resourceID))
}

// List returns audit logs for admin.
func (s *AuditService) List(q dto.AuditQuery) (*dto.PageData, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = util.DefaultPageSize
	}
	items, total, err := s.auditRepo.List(q)
	if err != nil {
		s.logger.Error("audit list failed", "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.AuditLogDTO, 0, len(items))
	for i := range items {
		list = append(list, dto.FromAuditLog(&items[i]))
	}
	return &dto.PageData{List: list, Total: total, Page: q.Page, Size: q.PageSize}, nil
}
