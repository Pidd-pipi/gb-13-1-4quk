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

// AuthHandler exposes auth endpoints.
type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger}
}

// SendCode sends an email verification code.
func (h *AuthHandler) SendCode(c *gin.Context) {
	var req dto.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "email 格式不正确: "+err.Error()))
		return
	}
	resp, err := h.authService.SendCode(c.Request.Context(), req.Email)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Register creates a student account.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "注册参数不合法: "+err.Error()))
		return
	}
	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Login authenticates a user.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "登录参数不合法: "+err.Error()))
		return
	}
	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}
