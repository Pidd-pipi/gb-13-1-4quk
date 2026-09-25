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

// WishHandler exposes wish endpoints.
type WishHandler struct {
	wishService *service.WishService
	logger      *slog.Logger
}

// NewWishHandler creates a WishHandler.
func NewWishHandler(wishService *service.WishService, logger *slog.Logger) *WishHandler {
	return &WishHandler{wishService: wishService, logger: logger}
}

// Create publishes a wish.
func (h *WishHandler) Create(c *gin.Context) {
	var req dto.CreateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "求购参数不合法: "+err.Error()))
		return
	}
	resp, err := h.wishService.CreateWish(middleware.GetUserID(c), req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// List searches wishes.
func (h *WishHandler) List(c *gin.Context) {
	var q dto.WishQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "查询参数不合法: "+err.Error()))
		return
	}
	resp, err := h.wishService.ListWishes(q)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Detail returns a single wish.
func (h *WishHandler) Detail(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "求购 id 不合法"))
		return
	}
	resp, err := h.wishService.GetWish(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Update edits a wish.
func (h *WishHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "求购 id 不合法"))
		return
	}
	var req dto.UpdateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "求购参数不合法: "+err.Error()))
		return
	}
	resp, err := h.wishService.UpdateWish(middleware.GetUserID(c), id, req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Close closes a wish.
func (h *WishHandler) Close(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "求购 id 不合法"))
		return
	}
	resp, err := h.wishService.CloseWish(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Delete removes a wish.
func (h *WishHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "求购 id 不合法"))
		return
	}
	if err := h.wishService.DeleteWish(middleware.GetUserID(c), id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"id": id}))
}
