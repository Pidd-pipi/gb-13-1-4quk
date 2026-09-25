package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/middleware"
	"github.com/campusbooks/campusbooks/internal/service"
	"github.com/campusbooks/campusbooks/internal/util"
)

// EvaluationHandler exposes evaluation endpoints.
type EvaluationHandler struct {
	evalService *service.EvaluationService
	logger      *slog.Logger
}

// NewEvaluationHandler creates an EvaluationHandler.
func NewEvaluationHandler(evalService *service.EvaluationService, logger *slog.Logger) *EvaluationHandler {
	return &EvaluationHandler{evalService: evalService, logger: logger}
}

// Create creates an evaluation for a finished trade.
func (h *EvaluationHandler) Create(c *gin.Context) {
	var req dto.CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "评价参数不合法: "+err.Error()))
		return
	}
	resp, err := h.evalService.CreateEvaluation(middleware.GetUserID(c), req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ListByUser lists evaluations received by a user.
func (h *EvaluationHandler) ListByUser(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "用户 id 不合法"))
		return
	}
	p := util.ParsePageParams(c)
	resp, err := h.evalService.ListEvaluations(id, p.Page, p.Size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}
