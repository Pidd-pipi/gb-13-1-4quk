package model

import (
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
)

// Borrow 短借申请：同学在书籍详情提交，卖家同意后进入借出状态，
// 借阅人归还、卖家确认收回后闭环。一本书同一时间至多存在一笔未结束的借阅。
type Borrow struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	BookID     uint `gorm:"index;not null" json:"book_id"`
	LenderID   uint `gorm:"index;not null" json:"lender_id"`   // 书籍卖家
	BorrowerID uint `gorm:"index;not null" json:"borrower_id"` // 借阅申请人
	// Duration 申请的借阅时长（天），取发布时书籍配置（7/14）。
	Duration int `gorm:"not null;default:0" json:"duration"`
	// Status pending / approved / rejected / returning / returned。
	Status string `gorm:"size:16;not null;default:pending;index" json:"status"`
	// DueAt 卖家同意时计算的到期日（approved+duration 天）。
	DueAt        *time.Time `gorm:"column:due_at;index" json:"due_at"`
	ApprovedAt   *time.Time `json:"approved_at"`
	ReturnedAt   *time.Time `json:"returned_at"`  // 借阅人发起归还的时间
	ConfirmedAt  *time.Time `json:"confirmed_at"` // 卖家确认收回的时间
	RejectReason string     `gorm:"size:255" json:"reject_reason"`
	RemindedAt   *time.Time `json:"reminded_at"` // 卖家最近一次发送逾期提醒的时间
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Book     *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Lender   *User `gorm:"foreignKey:LenderID" json:"lender,omitempty"`
	Borrower *User `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
}

// IsActive 申请是否尚未终结（approved/returning 仍在借阅周期内）。
func (b *Borrow) IsActive() bool {
	return b.Status == constants.BorrowStatusApproved || b.Status == constants.BorrowStatusReturning
}

// IsPending 是否待卖家处理。
func (b *Borrow) IsPending() bool { return b.Status == constants.BorrowStatusPending }

// IsOverdue 是否已超过到期日仍未归还（含待确认归还阶段）。
func (b *Borrow) IsOverdue(now time.Time) bool {
	return b.IsActive() && b.DueAt != nil && now.After(*b.DueAt)
}
