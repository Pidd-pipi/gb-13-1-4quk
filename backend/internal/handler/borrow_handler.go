package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/middleware"
	"github.com/campusbooks/campusbooks/internal/service"
	"github.com/campusbooks/campusbooks/internal/util"
)

// BorrowHandler exposes short-borrow endpoints.
type BorrowHandler struct {
	borrowService *service.BorrowService
	logger        *slog.Logger
}

// NewBorrowHandler creates a BorrowHandler.
func NewBorrowHandler(borrowService *service.BorrowService, logger *slog.Logger) *BorrowHandler {
	return &BorrowHandler{borrowService: borrowService, logger: logger}
}

// Create submits a borrow request for a book.
func (h *BorrowHandler) Create(c *gin.Context) {
	bookID, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	var req dto.CreateBorrowRequest
	// 借阅附言可空：忽略空 body 的 EOF，仅拒绝结构不合法的 JSON
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请参数不合法: "+err.Error()))
		return
	}
	resp, err := h.borrowService.Apply(middleware.GetUserID(c), bookID, req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ListByBook returns borrow requests of a book (seller sees all, borrower sees own).
func (h *BorrowHandler) ListByBook(c *gin.Context) {
	bookID, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	q := dto.BorrowQuery{BookID: bookID}
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "查询参数不合法: "+err.Error()))
		return
	}
	if q.Role == "" {
		q.Role = "lender"
	}
	fillBorrowScope(c, &q)
	resp, err := h.borrowService.List(q)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// List returns the current user's borrow requests (role=lender/borrower).
func (h *BorrowHandler) List(c *gin.Context) {
	var q dto.BorrowQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "查询参数不合法: "+err.Error()))
		return
	}
	if q.Role == "" {
		q.Role = "borrower"
	}
	fillBorrowScope(c, &q)
	resp, err := h.borrowService.List(q)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Detail returns one borrow request.
func (h *BorrowHandler) Detail(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	resp, err := h.borrowService.Get(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Approve lets the seller approve a pending request.
func (h *BorrowHandler) Approve(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	resp, err := h.borrowService.Approve(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Reject lets the seller reject a pending request.
func (h *BorrowHandler) Reject(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	var req dto.RejectBorrowRequest
	_ = c.ShouldBindJSON(&req) // 拒绝原因可空
	resp, err := h.borrowService.Reject(middleware.GetUserID(c), id, req)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Return lets the borrower mark the book as returned (awaiting confirmation).
func (h *BorrowHandler) Return(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	resp, err := h.borrowService.Return(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// ConfirmReturn lets the seller confirm the book is back and lendable again.
func (h *BorrowHandler) ConfirmReturn(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	resp, err := h.borrowService.ConfirmReturn(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Remind lets the seller send an overdue reminder to the borrower.
func (h *BorrowHandler) Remind(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	var req dto.RemindBorrowRequest
	_ = c.ShouldBindJSON(&req) // 附言可空，使用默认提醒文案
	resp, err := h.borrowService.Remind(middleware.GetUserID(c), id, req)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Contact opens the lender/borrower conversation for a borrow.
func (h *BorrowHandler) Contact(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "借阅申请 id 不合法"))
		return
	}
	var req dto.CreateConversationRequest
	_ = c.ShouldBindJSON(&req) // 首条消息可空
	resp, err := h.borrowService.Contact(middleware.GetUserID(c), id, req.Content)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// fillBorrowScope constrains the query to the current user by role.
func fillBorrowScope(c *gin.Context, q *dto.BorrowQuery) {
	uid := middleware.GetUserID(c)
	if q.Role == "lender" {
		q.LenderID = uid
	} else {
		q.BorrowerID = uid
	}
}
