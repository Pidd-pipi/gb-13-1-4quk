package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerWishRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.WishHandler, convHandler *handler.ConversationHandler, limiter *middleware.RateLimiter) {
	wishes := g.Group("/wishes")
	{
		wishes.GET("", h.List)
		wishes.POST("", middleware.AuthRequired(cfg), h.Create)
		wishes.GET("/:id", h.Detail)
		wishes.PUT("/:id", middleware.AuthRequired(cfg), h.Update)
		wishes.DELETE("/:id", middleware.AuthRequired(cfg), h.Delete)
		wishes.POST("/:id/close", middleware.AuthRequired(cfg), h.Close)
		wishes.POST("/:id/contact", middleware.AuthRequired(cfg), convHandler.Create)
	}
}
