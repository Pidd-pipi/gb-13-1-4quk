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

// BookHandler exposes book endpoints.
type BookHandler struct {
	bookService *service.BookService
	logger      *slog.Logger
}

// NewBookHandler creates a BookHandler.
func NewBookHandler(bookService *service.BookService, logger *slog.Logger) *BookHandler {
	return &BookHandler{bookService: bookService, logger: logger}
}

// Create publishes a book.
func (h *BookHandler) Create(c *gin.Context) {
	var req dto.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍参数不合法: "+err.Error()))
		return
	}
	resp, err := h.bookService.CreateBook(middleware.GetUserID(c), req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// List searches books.
func (h *BookHandler) List(c *gin.Context) {
	var q dto.BookQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "查询参数不合法: "+err.Error()))
		return
	}
	resp, err := h.bookService.ListBooks(q)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Detail returns a book and records browse history.
func (h *BookHandler) Detail(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	userID := middleware.GetUserID(c)
	resp, err := h.bookService.GetBookDetail(userID, id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Update edits a book.
func (h *BookHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	var req dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍参数不合法: "+err.Error()))
		return
	}
	resp, err := h.bookService.UpdateBook(middleware.GetUserID(c), id, req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Delete removes a book.
func (h *BookHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	if err := h.bookService.DeleteBook(middleware.GetUserID(c), id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"id": id}))
}

// Reserve marks a book as reserved.
func (h *BookHandler) Reserve(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	resp, err := h.bookService.ReserveBook(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// CancelReserve cancels a reservation.
func (h *BookHandler) CancelReserve(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	resp, err := h.bookService.CancelReserve(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Sold marks a book as sold.
func (h *BookHandler) Sold(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	resp, err := h.bookService.MarkSold(middleware.GetUserID(c), id)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// AddFavorite favorites a book.
func (h *BookHandler) AddFavorite(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	if err := h.bookService.AddFavorite(middleware.GetUserID(c), id); err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"id": id, "favorited": true}))
}

// RemoveFavorite un-favorites a book.
func (h *BookHandler) RemoveFavorite(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "书籍 id 不合法"))
		return
	}
	if err := h.bookService.RemoveFavorite(middleware.GetUserID(c), id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(gin.H{"id": id, "favorited": false}))
}

// Favorites lists the user's favorites.
func (h *BookHandler) Favorites(c *gin.Context) {
	p := util.ParsePageParams(c)
	resp, err := h.bookService.ListFavorites(middleware.GetUserID(c), p.Page, p.Size)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// History lists the user's browse history.
func (h *BookHandler) History(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	resp, err := h.bookService.ListHistory(middleware.GetUserID(c), limit)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

// Recommendations lists same-department books.
func (h *BookHandler) Recommendations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := h.bookService.GetRecommendations(middleware.GetUserID(c), limit)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(resp))
}

func parseID(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
