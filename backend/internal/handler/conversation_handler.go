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

// ConversationHandler exposes conversation/message endpoints.
type ConversationHandler struct {
	convService *service.ConversationService
	logger      *slog.Logger
}

// NewConversationHandler creates a ConversationHandler.
func NewConversationHandler(convService *service.ConversationService, logger *slog.Logger) *ConversationHandler {
	return &ConversationHandler{convService: convService, logger: logger}
}

// Create opens a conversation from a book or wish.
func (h *ConversationHandler) Create(c *gin.Context) {
	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "会话参数不合法: "+err.Error()))
		return
	}
	if req.BookID == 0 && req.WishID == 0 {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "book_id 与 wish_id 至少提供一个"))
		return
	}
	userID := middleware.GetUserID(c)
	var resp interface{}
	var err error
	switch {
	case req.BookID > 0 && req.ToUserID > 0:
		resp, err = h.convService.CreateConversationToUser(userID, req.BookID, req.ToUserID, req.Content)
	case req.BookID > 0:
		resp, err = h.convService.CreateConversationFromBook(userID, req.BookID, req.Content)
	default:
		resp, err = h.convService.CreateConversationFromWish(userID, req.WishID, req.Content)
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// List lists the current user's conversations.
func (h *ConversationHandler) List(c *gin.Context) {
	p := util.ParsePageParams(c)
	resp, err := h.convService.ListConversations(middleware.GetUserID(c), p.Page, p.Size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Detail returns a conversation.
func (h *ConversationHandler) Detail(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "会话 id 不合法"))
		return
	}
	resp, err := h.convService.GetConversation(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// SendMessage sends a message in a conversation.
func (h *ConversationHandler) SendMessage(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "会话 id 不合法"))
		return
	}
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "消息参数不合法: "+err.Error()))
		return
	}
	resp, err := h.convService.SendMessage(id, middleware.GetUserID(c), req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ListMessages lists messages of a conversation.
func (h *ConversationHandler) ListMessages(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "会话 id 不合法"))
		return
	}
	resp, err := h.convService.ListMessages(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// MarkRead marks a conversation read.
func (h *ConversationHandler) MarkRead(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "会话 id 不合法"))
		return
	}
	if err := h.convService.MarkRead(middleware.GetUserID(c), id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"id": id}))
}
