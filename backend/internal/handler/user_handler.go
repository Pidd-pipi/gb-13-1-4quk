package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/middleware"
	"github.com/campusbooks/campusbooks/internal/service"
	"github.com/campusbooks/campusbooks/internal/util"
)

// UserHandler exposes user profile/management endpoints.
type UserHandler struct {
	userService *service.UserService
	logger      *slog.Logger
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// Me returns the current user profile.
func (h *UserHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	resp, err := h.userService.GetProfile(userID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// UpdateProfile updates the current user profile.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "资料参数不合法: "+err.Error()))
		return
	}
	resp, err := h.userService.UpdateProfile(middleware.GetUserID(c), req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// UpdateAvatar updates the current user avatar.
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	var req dto.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "头像参数不合法: "+err.Error()))
		return
	}
	resp, err := h.userService.UpdateAvatar(middleware.GetUserID(c), req.AvatarURL)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ListUsers lists users (admin).
func (h *UserHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	p := util.ParsePageParams(c)
	resp, err := h.userService.ListUsers(role, p.Page, p.Size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// GetUser returns a user profile by id.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "用户 id 不合法"))
		return
	}
	resp, err := h.userService.GetProfile(uint(id))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// GetStats returns user homepage stats.
func (h *UserHandler) GetStats(c *gin.Context) {
	var userID uint
	if v := c.Param("id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "用户 id 不合法"))
			return
		}
		userID = uint(id)
	} else {
		userID = middleware.GetUserID(c)
	}
	resp, err := h.userService.GetStats(userID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}
