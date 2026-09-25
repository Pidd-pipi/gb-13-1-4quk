package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/service"
)

// AuditMiddleware records state-changing requests into the audit log.
func AuditMiddleware(auditService *service.AuditService, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
			return
		}
		userID := GetUserID(c)
		auditService.Record(
			userID,
			method+" "+c.FullPath(),
			firstSegment(c.FullPath()),
			0,
			truncateDetail(c.Request.URL.RequestURI()),
			c.ClientIP(),
			GetRequestID(c),
		)
	}
}

func firstSegment(path string) string {
	for i := 1; i < len(path); i++ {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return path
}

func truncateDetail(s string) string {
	if len(s) > 200 {
		return s[:200]
	}
	return s
}
