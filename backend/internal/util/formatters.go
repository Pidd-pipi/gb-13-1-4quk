package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
)

// 格式化工具集中管理：日期、状态文本、类型文本、价格文本等。
// 新增枚举/状态时必须同步修改这里与 constants/enums.go、前端 constants/enums.ts。

// FormatTime renders time as "2006-01-02 15:04".
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04")
}

// FormatBookStatusText maps a book status to Chinese text.
func FormatBookStatusText(status string) string {
	switch status {
	case constants.BookStatusOnSale:
		return "在售"
	case constants.BookStatusReserved:
		return "已预约"
	case constants.BookStatusSold:
		return "已售出"
	case constants.BookStatusLoaned:
		return "借出中"
	default:
		return "未知"
	}
}

// FormatBorrowStatusText maps a borrow request status to Chinese text.
func FormatBorrowStatusText(status string) string {
	switch status {
	case constants.BorrowStatusPending:
		return "待同意"
	case constants.BorrowStatusApproved:
		return "借出中"
	case constants.BorrowStatusRejected:
		return "已拒绝"
	case constants.BorrowStatusReturning:
		return "待确认归还"
	case constants.BorrowStatusReturned:
		return "已归还"
	default:
		return "未知"
	}
}

// FormatConditionText maps a condition code to Chinese text.
func FormatConditionText(cond string) string {
	switch cond {
	case constants.ConditionBrandNew:
		return "全新"
	case constants.ConditionNineNew:
		return "九成新"
	case constants.ConditionSevenNew:
		return "七成新"
	case constants.ConditionFiveNew:
		return "五成新"
	default:
		return "未知"
	}
}

// FormatSubjectText maps a subject category code to Chinese text.
func FormatSubjectText(subject string) string {
	switch subject {
	case constants.SubjectScience:
		return "理工"
	case constants.SubjectHumanities:
		return "文史"
	case constants.SubjectEconManagement:
		return "经管"
	case constants.SubjectArt:
		return "艺术"
	case constants.SubjectOther:
		return "其他"
	default:
		return "未知"
	}
}

// FormatTradeTypeText maps a trade type code to Chinese text.
func FormatTradeTypeText(tradeType string) string {
	switch tradeType {
	case constants.TradeTypeInPerson:
		return "面交"
	case constants.TradeTypeMail:
		return "邮寄"
	default:
		return "未知"
	}
}

// FormatEvaluationTypeText maps an evaluation type code to Chinese text.
func FormatEvaluationTypeText(t string) string {
	switch t {
	case constants.EvaluationGood:
		return "好评"
	case constants.EvaluationNeutral:
		return "中评"
	case constants.EvaluationBad:
		return "差评"
	default:
		return "未知"
	}
}

// FormatRoleText maps a role code to Chinese text.
func FormatRoleText(role string) string {
	switch role {
	case constants.RoleAdmin:
		return "管理员"
	case constants.RoleStudent:
		return "学生"
	default:
		return "未知"
	}
}

// FormatWishStatusText maps a wish status code to Chinese text.
func FormatWishStatusText(status string) string {
	switch status {
	case constants.WishStatusOpen:
		return "进行中"
	case constants.WishStatusClosed:
		return "已关闭"
	default:
		return "未知"
	}
}

// FormatPrice renders a float64 price as "¥12.50".
func FormatPrice(price float64) string {
	return fmt.Sprintf("¥%.2f", price)
}

// JoinImages joins image urls with comma for JSON column storage.
func JoinImages(images []string) string { return strings.Join(images, ",") }

// SplitImages splits the comma joined image column into a slice.
func SplitImages(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
