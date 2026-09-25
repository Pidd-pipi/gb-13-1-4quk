package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerBookRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.BookHandler, limiter *middleware.RateLimiter) {
	books := g.Group("/books")
	{
		books.GET("", h.List)
		books.GET("/recommendations", middleware.AuthRequired(cfg), h.Recommendations)
		books.GET("/favorites", middleware.AuthRequired(cfg), h.Favorites)
		books.GET("/history", middleware.AuthRequired(cfg), h.History)
		books.POST("", middleware.AuthRequired(cfg), h.Create)
		books.GET("/:id", h.Detail)
		books.PUT("/:id", middleware.AuthRequired(cfg), h.Update)
		books.DELETE("/:id", middleware.AuthRequired(cfg), h.Delete)
		books.POST("/:id/reserve", middleware.AuthRequired(cfg), h.Reserve)
		books.POST("/:id/cancel-reserve", middleware.AuthRequired(cfg), h.CancelReserve)
		books.POST("/:id/sold", middleware.AuthRequired(cfg), h.Sold)
		books.POST("/:id/favorite", middleware.AuthRequired(cfg), h.AddFavorite)
		books.DELETE("/:id/favorite", middleware.AuthRequired(cfg), h.RemoveFavorite)
	}
}
