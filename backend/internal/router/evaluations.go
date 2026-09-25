package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerEvaluationRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.EvaluationHandler, limiter *middleware.RateLimiter) {
	evals := g.Group("/evaluations")
	{
		evals.POST("", middleware.AuthRequired(cfg), h.Create)
	}
}
