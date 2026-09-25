package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
)

// RequestLogger logs every request with request_id/method/path/status/latency.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info(constants.LogRequestHandled,
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
