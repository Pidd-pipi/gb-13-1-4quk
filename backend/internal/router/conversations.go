package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerConversationRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.ConversationHandler, limiter *middleware.RateLimiter) {
	convs := g.Group("/conversations")
	convs.Use(middleware.AuthRequired(cfg))
	{
		convs.GET("", h.List)
		convs.POST("", h.Create)
		convs.GET("/:id", h.Detail)
		convs.POST("/:id/messages", h.SendMessage)
		convs.GET("/:id/messages", h.ListMessages)
		convs.PUT("/:id/read", h.MarkRead)
	}
}
