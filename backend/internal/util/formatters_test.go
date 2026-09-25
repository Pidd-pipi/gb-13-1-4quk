package util

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
)

func TestFormatBookStatusText(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{constants.BookStatusOnSale, "在售"},
		{constants.BookStatusReserved, "已预约"},
		{constants.BookStatusSold, "已售出"},
		{"unknown", "未知"},
	}
	for _, c := range cases {
		if got := FormatBookStatusText(c.in); got != c.want {
			t.Errorf("FormatBookStatusText(%q)=%q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatConditionText(t *testing.T) {
	if got := FormatConditionText(constants.ConditionBrandNew); got != "全新" {
		t.Errorf("brand_new -> %q, want 全新", got)
	}
	if got := FormatConditionText("x"); got != "未知" {
		t.Errorf("unknown -> %q, want 未知", got)
	}
}

func TestSplitJoinImages(t *testing.T) {
	joined := JoinImages([]string{"a.jpg", "b.png"})
	if joined != "a.jpg,b.png" {
		t.Errorf("join = %q", joined)
	}
	parts := SplitImages(joined)
	if len(parts) != 2 || parts[0] != "a.jpg" {
		t.Errorf("split = %v", parts)
	}
	if got := SplitImages(""); len(got) != 0 {
		t.Errorf("empty split = %v", got)
	}
}

func TestFormatPrice(t *testing.T) {
	if got := FormatPrice(12.5); got != "¥12.50" {
		t.Errorf("FormatPrice(12.5)=%q", got)
	}
}
