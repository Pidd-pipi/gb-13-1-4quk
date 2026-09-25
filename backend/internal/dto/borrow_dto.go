package dto

import (
	"time"

	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// CreateBorrowRequest 提交借阅申请请求（时长以书籍发布配置为准，无需前端传）。
type CreateBorrowRequest struct {
	Message string `json:"message" binding:"max=200"`
}

// RejectBorrowRequest 卖家拒绝借阅申请。
type RejectBorrowRequest struct {
	Reason string `json:"reason" binding:"max=255"`
}

// RemindBorrowRequest 卖家发送还书提醒（可附言，为空时使用默认提醒文案）。
type RemindBorrowRequest struct {
	Content string `json:"content" binding:"max=500"`
}

// BorrowQuery 借阅申请列表查询：LenderID/BorrowerID 由 handler 按当前登录用户
// 与 role 参数填充，可按书籍与状态过滤。
type BorrowQuery struct {
	Role       string `form:"role"`
	LenderID   uint   `form:"-"`
	BorrowerID uint   `form:"-"`
	BookID     uint   `form:"book_id"`
	Status     string `form:"status"`
	Page       int    `form:"page"`
	Size       int    `form:"page_size"`
}

// BorrowDTO 短借申请响应体。
type BorrowDTO struct {
	ID           uint     `json:"id"`
	BookID       uint     `json:"book_id"`
	LenderID     uint     `json:"lender_id"`
	BorrowerID   uint     `json:"borrower_id"`
	Duration     int      `json:"duration"`
	Status       string   `json:"status"`
	StatusText   string   `json:"status_text"`
	DueAt        string   `json:"due_at"`
	IsOverdue    bool     `json:"is_overdue"`
	ApprovedAt   string   `json:"approved_at"`
	ReturnedAt   string   `json:"returned_at"`
	ConfirmedAt  string   `json:"confirmed_at"`
	RejectReason string   `json:"reject_reason"`
	RemindedAt   string   `json:"reminded_at"`
	CreatedAt    string   `json:"created_at"`
	Book         *BookDTO `json:"book,omitempty"`
	Lender       *UserDTO `json:"lender,omitempty"`
	Borrower     *UserDTO `json:"borrower,omitempty"`
}

// FromBorrow converts a model.Borrow to BorrowDTO. Overdue is computed against
// now so expired loans surface 逾期 without a background job.
func FromBorrow(b *model.Borrow, now time.Time) BorrowDTO {
	d := BorrowDTO{
		ID:           b.ID,
		BookID:       b.BookID,
		LenderID:     b.LenderID,
		BorrowerID:   b.BorrowerID,
		Duration:     b.Duration,
		Status:       b.Status,
		StatusText:   util.FormatBorrowStatusText(b.Status),
		IsOverdue:    b.IsOverdue(now),
		RejectReason: b.RejectReason,
		CreatedAt:    util.FormatTime(b.CreatedAt),
	}
	if b.DueAt != nil {
		d.DueAt = util.FormatTime(*b.DueAt)
	}
	if b.ApprovedAt != nil {
		d.ApprovedAt = util.FormatTime(*b.ApprovedAt)
	}
	if b.ReturnedAt != nil {
		d.ReturnedAt = util.FormatTime(*b.ReturnedAt)
	}
	if b.ConfirmedAt != nil {
		d.ConfirmedAt = util.FormatTime(*b.ConfirmedAt)
	}
	if b.RemindedAt != nil {
		d.RemindedAt = util.FormatTime(*b.RemindedAt)
	}
	return d
}
