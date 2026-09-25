package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// BorrowDTO 借阅申请响应体。
type BorrowDTO struct {
	ID          uint     `json:"id"`
	BookID      uint     `json:"book_id"`
	BorrowerID  uint     `json:"borrower_id"`
	SellerID    uint     `json:"seller_id"`
	Status      string   `json:"status"`
	StatusText  string   `json:"status_text"`
	ApprovedAt  string   `json:"approved_at"`
	DueAt       string   `json:"due_at"`
	ReturnedAt  string   `json:"returned_at"`
	ConfirmedAt string   `json:"confirmed_at"`
	RemindedAt  string   `json:"reminded_at"`
	Overdue     bool     `json:"overdue"`
	CreatedAt   string   `json:"created_at"`
	Book        *BookDTO `json:"book,omitempty"`
	Borrower    *UserDTO `json:"borrower,omitempty"`
	Seller      *UserDTO `json:"seller,omitempty"`
}

// FromBorrow converts a model.BorrowRequest to BorrowDTO.
func FromBorrow(r *model.BorrowRequest) BorrowDTO {
	d := BorrowDTO{
		ID:         r.ID,
		BookID:     r.BookID,
		BorrowerID: r.BorrowerID,
		SellerID:   r.SellerID,
		Status:     r.Status,
		StatusText: util.FormatBorrowStatusText(r.Status),
		Overdue:    r.IsOverdue(),
		CreatedAt:  util.FormatTime(r.CreatedAt),
	}
	if r.ApprovedAt != nil {
		d.ApprovedAt = util.FormatTime(*r.ApprovedAt)
	}
	if r.DueAt != nil {
		d.DueAt = util.FormatTime(*r.DueAt)
	}
	if r.ReturnedAt != nil {
		d.ReturnedAt = util.FormatTime(*r.ReturnedAt)
	}
	if r.ConfirmedAt != nil {
		d.ConfirmedAt = util.FormatTime(*r.ConfirmedAt)
	}
	if r.RemindedAt != nil {
		d.RemindedAt = util.FormatTime(*r.RemindedAt)
	}
	if r.Book != nil {
		b := FromBook(r.Book)
		d.Book = &b
	}
	if r.Borrower != nil {
		u := FromUser(r.Borrower)
		d.Borrower = &u
	}
	if r.Seller != nil {
		u := FromUser(r.Seller)
		d.Seller = &u
	}
	return d
}
