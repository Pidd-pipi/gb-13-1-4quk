package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerBorrowRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.BorrowHandler, limiter *middleware.RateLimiter) {
	books := g.Group("/books")
	{
		books.POST("/:id/borrow-requests", middleware.AuthRequired(cfg), h.Apply)
		books.GET("/:id/borrow-requests", middleware.AuthRequired(cfg), h.ListByBook)
	}
	borrows := g.Group("/borrow-requests")
	{
		borrows.GET("", middleware.AuthRequired(cfg), h.ListMine)
		borrows.POST("/:id/approve", middleware.AuthRequired(cfg), h.Approve)
		borrows.POST("/:id/reject", middleware.AuthRequired(cfg), h.Reject)
		borrows.POST("/:id/return", middleware.AuthRequired(cfg), h.Return)
		borrows.POST("/:id/confirm-return", middleware.AuthRequired(cfg), h.ConfirmReturn)
		borrows.POST("/:id/remind", middleware.AuthRequired(cfg), h.Remind)
	}
}
