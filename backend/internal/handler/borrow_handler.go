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

// BorrowHandler exposes short-term lending endpoints.
type BorrowHandler struct {
	borrowService *service.BorrowService
	logger        *slog.Logger
}

// NewBorrowHandler creates a BorrowHandler.
func NewBorrowHandler(borrowService *service.BorrowService, logger *slog.Logger) *BorrowHandler {
	return &BorrowHandler{borrowService: borrowService, logger: logger}
}

// Apply submits a borrow request for a book.
func (h *BorrowHandler) Apply(c *gin.Context) {
	bookID, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	resp, err := h.borrowService.Apply(middleware.GetUserID(c), bookID)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ListByBook lists borrow requests of a book (seller only).
func (h *BorrowHandler) ListByBook(c *gin.Context) {
	bookID, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	resp, err := h.borrowService.ListByBook(middleware.GetUserID(c), bookID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ListMine lists the caller's borrow requests (role=borrower|seller).
func (h *BorrowHandler) ListMine(c *gin.Context) {
	role := c.DefaultQuery("role", "borrower")
	if role != "borrower" && role != "seller" {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "role 仅支持 borrower / seller"))
		return
	}
	p := util.ParsePageParams(c)
	resp, err := h.borrowService.ListMine(middleware.GetUserID(c), role, p.Page, p.Size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Approve lets the seller approve a borrow request.
func (h *BorrowHandler) Approve(c *gin.Context) {
	h.transition(c, h.borrowService.Approve)
}

// Reject lets the seller reject a borrow request.
func (h *BorrowHandler) Reject(c *gin.Context) {
	h.transition(c, h.borrowService.Reject)
}

// Return lets the borrower mark the book as returned.
func (h *BorrowHandler) Return(c *gin.Context) {
	h.transition(c, h.borrowService.Return)
}

// ConfirmReturn lets the seller confirm the book is back.
func (h *BorrowHandler) ConfirmReturn(c *gin.Context) {
	h.transition(c, h.borrowService.ConfirmReturn)
}

// Remind lets the seller send a return reminder to the borrower.
func (h *BorrowHandler) Remind(c *gin.Context) {
	h.transition(c, h.borrowService.Remind)
}

func (h *BorrowHandler) transition(c *gin.Context, fn func(userID, requestID uint) (*dto.BorrowDTO, error)) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	resp, err := fn(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}
