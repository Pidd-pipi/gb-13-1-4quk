package router

import (
	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
)

func registerUserRoutes(g *gin.RouterGroup, cfg *config.Config, h *handler.UserHandler, evalHandler *handler.EvaluationHandler, limiter *middleware.RateLimiter) {
	users := g.Group("/users")
	{
		users.GET("", middleware.AuthRequired(cfg), middleware.RequireRole(constants.RoleAdmin), h.ListUsers)
		users.GET("/me", middleware.AuthRequired(cfg), h.Me)
		users.PUT("/me", middleware.AuthRequired(cfg), h.UpdateProfile)
		users.PUT("/me/avatar", middleware.AuthRequired(cfg), h.UpdateAvatar)
		users.GET("/me/stats", middleware.AuthRequired(cfg), h.GetStats)
		users.GET("/:id", middleware.AuthRequired(cfg), h.GetUser)
		users.GET("/:id/stats", h.GetStats)
		users.GET("/:id/evaluations", evalHandler.ListByUser)
	}
}
