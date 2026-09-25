package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerUploadRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.UploadHandler, limiter *middleware.RateLimiter) {
	g.GET("/files/:key", h.Get)
	uploads := g.Group("/uploads")
	uploads.Use(middleware.AuthRequired(cfg))
	{
		uploads.POST("", h.Upload)
	}
}
