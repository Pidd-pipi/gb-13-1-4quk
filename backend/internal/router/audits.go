package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerAuditRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.AuditHandler, limiter *middleware.RateLimiter) {
	audits := g.Group("/audit-logs")
	audits.Use(middleware.AuthRequired(cfg), middleware.RequireRole(constants.RoleAdmin))
	{
		audits.GET("", h.List)
	}
}
