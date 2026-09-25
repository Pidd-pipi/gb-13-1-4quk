package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
)

type bucket struct {
	count   int
	resetAt time.Time
}

// RateLimiter is a simple in-memory sliding window limiter per IP.
type RateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	items  map[string]*bucket
	logger *slog.Logger
}

// NewRateLimiter creates a RateLimiter.
func NewRateLimiter(limit int, window time.Duration, logger *slog.Logger) *RateLimiter {
	return &RateLimiter{limit: limit, window: window, items: make(map[string]*bucket), logger: logger}
}

// Limit returns a gin middleware that rejects bursty requests.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		rl.mu.Lock()
		b, ok := rl.items[ip]
		if !ok || now.After(b.resetAt) {
			b = &bucket{count: 0, resetAt: now.Add(rl.window)}
			rl.items[ip] = b
		}
		b.count++
		over := b.count > rl.limit
		rl.mu.Unlock()
		if over {
			rl.logger.Warn(constants.LogRateLimited, "path", c.Request.URL.Path, "ip", ip)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.Fail(constants.CodeRateLimited, constants.MsgRateLimited))
			return
		}
		c.Next()
	}
}
