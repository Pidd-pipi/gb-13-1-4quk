package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
)

// Recovery catches panics and returns a unified 500 response.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(constants.LogPanicRecovered, "request_id", GetRequestID(c), "error", r, "stack", string(debug.Stack()))
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Fail(constants.CodeInternalError, constants.MsgInternalError))
			}
		}()
		c.Next()
	}
}
