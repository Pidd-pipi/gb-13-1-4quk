package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerBorrowRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.BorrowHandler) {
	// 我的借阅（/borrows 需在 /books/:id/... 之外独立分组）
	borrows := g.Group("/borrows")
	{
		borrows.GET("", middleware.AuthRequired(cfg), h.List)
		borrows.GET("/:id", middleware.AuthRequired(cfg), h.Detail)
		borrows.POST("/:id/approve", middleware.AuthRequired(cfg), h.Approve)
		borrows.POST("/:id/reject", middleware.AuthRequired(cfg), h.Reject)
		borrows.POST("/:id/return", middleware.AuthRequired(cfg), h.Return)
		borrows.POST("/:id/confirm-return", middleware.AuthRequired(cfg), h.ConfirmReturn)
		borrows.POST("/:id/remind", middleware.AuthRequired(cfg), h.Remind)
		borrows.POST("/:id/contact", middleware.AuthRequired(cfg), h.Contact)
	}

	// 某本书的借阅申请列表（挂在书籍资源下）
	books := g.Group("/books")
	{
		books.GET("/:id/borrows", middleware.AuthRequired(cfg), h.ListByBook)
		books.POST("/:id/borrows", middleware.AuthRequired(cfg), h.Create)
	}
}
