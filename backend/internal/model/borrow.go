package model

import (
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
)

// BorrowRequest 短借申请：同学对可借书籍发起，卖家同意后进入借出，
// 借阅人归还、卖家确认拿回后完成。
type BorrowRequest struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	BookID      uint       `gorm:"index;not null" json:"book_id"`
	BorrowerID  uint       `gorm:"index;not null" json:"borrower_id"`
	SellerID    uint       `gorm:"index;not null" json:"seller_id"`
	Status      string     `gorm:"size:16;default:pending;index" json:"status"`
	ApprovedAt  *time.Time `json:"approved_at"`
	DueAt       *time.Time `json:"due_at"`
	ReturnedAt  *time.Time `json:"returned_at"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	RemindedAt  *time.Time `json:"reminded_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Book     *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Borrower *User `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
	Seller   *User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}

// IsLentOut reports whether the request is in an active lent-out state
// (approved or returned-but-unconfirmed), i.e. the book is still out.
func (r *BorrowRequest) IsLentOut() bool {
	return r.Status == constants.BorrowStatusApproved || r.Status == constants.BorrowStatusReturned
}

// IsOverdue reports whether the lent-out book has passed its due date.
func (r *BorrowRequest) IsOverdue() bool {
	return r.IsLentOut() && r.DueAt != nil && time.Now().After(*r.DueAt)
}
