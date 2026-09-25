package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/util"
)

// ErrorHandler converts panics/app errors into the unified response.
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			c.AbortWithStatusJSON(appErr.HTTPStatus, dto.Fail(appErr.Code, appErr.Message))
			return
		}
		logger.Error("unhandled error", "error", err, "request_id", GetRequestID(c), "path", c.Request.URL.Path)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Fail(constants.CodeInternalError, constants.MsgInternalError))
	}
}
