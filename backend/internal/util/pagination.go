package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// PageParams holds normalized pagination parameters.
type PageParams struct {
	Page int
	Size int
}

// ParsePageParams reads page/page_size from query with defaults.
func ParsePageParams(c *gin.Context) PageParams {
	page := parsePositiveInt(c.Query("page"), 1)
	size := parsePositiveInt(c.Query("page_size"), DefaultPageSize)
	if size > MaxPageSize {
		size = MaxPageSize
	}
	return PageParams{Page: page, Size: size}
}

// Offset returns the SQL offset for the pagination params.
func (p PageParams) Offset() int {
	return (p.Page - 1) * p.Size
}

func parsePositiveInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
