package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerAuthRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.AuthHandler, limiter *middleware.RateLimiter) {
	auth := g.Group("/auth")
	{
		auth.POST("/send-code", limiter.Limit(), h.SendCode)
		auth.POST("/register", limiter.Limit(), h.Register)
		auth.POST("/login", limiter.Limit(), h.Login)
	}
}
