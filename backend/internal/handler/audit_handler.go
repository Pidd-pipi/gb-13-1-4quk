package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/service"
	"github.com/campusbooks/campusbooks/internal/util"
)

// AuditHandler exposes admin audit endpoints.
type AuditHandler struct {
	auditService *service.AuditService
	logger       *slog.Logger
}

// NewAuditHandler creates an AuditHandler.
func NewAuditHandler(auditService *service.AuditService, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{auditService: auditService, logger: logger}
}

// List returns audit logs (admin only).
func (h *AuditHandler) List(c *gin.Context) {
	var q dto.AuditQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "查询参数不合法: "+err.Error()))
		return
	}
	resp, err := h.auditService.List(q)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}
