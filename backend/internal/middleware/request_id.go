package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID injects a request id into the context and response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(RequestIDKey, rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}
