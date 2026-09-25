package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/model"
)

// BorrowRepository handles borrow request persistence.
type BorrowRepository struct{ db *gorm.DB }

// NewBorrowRepository creates a BorrowRepository.
func NewBorrowRepository(db *gorm.DB) *BorrowRepository { return &BorrowRepository{db: db} }

// Create inserts a borrow request.
func (r *BorrowRepository) Create(req *model.BorrowRequest) error {
	return translate(r.db.Create(req).Error)
}

// FindByID locates a borrow request by id with relations preloaded.
func (r *BorrowRepository) FindByID(id uint) (*model.BorrowRequest, error) {
	var req model.BorrowRequest
	err := translate(r.db.Preload("Book").Preload("Borrower").Preload("Seller").First(&req, id).Error)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// FindByIDForUpdate locks a borrow request row for status transitions.
func (r *BorrowRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.BorrowRequest, error) {
	var req model.BorrowRequest
	err := translate(tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&req, id).Error)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// FindByIDTx returns a borrow request inside a transaction with relations preloaded.
func (r *BorrowRepository) FindByIDTx(tx *gorm.DB, id uint) (*model.BorrowRequest, error) {
	var req model.BorrowRequest
	err := translate(tx.Preload("Book").Preload("Borrower").Preload("Seller").First(&req, id).Error)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// Update persists a borrow request.
func (r *BorrowRepository) Update(req *model.BorrowRequest) error {
	return translate(r.db.Save(req).Error)
}

// UpdateTx persists a borrow request inside a transaction.
func (r *BorrowRepository) UpdateTx(tx *gorm.DB, req *model.BorrowRequest) error {
	return translate(tx.Save(req).Error)
}

// FindActiveByBook returns the lent-out (approved/returned) request for a book, if any.
func (r *BorrowRepository) FindActiveByBook(bookID uint) (*model.BorrowRequest, error) {
	var req model.BorrowRequest
	err := translate(r.db.Preload("Borrower").
		Where("book_id = ? AND status IN ?", bookID, []string{constants.BorrowStatusApproved, constants.BorrowStatusReturned}).
		Order("id DESC").First(&req).Error)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// FindOngoingByBookAndBorrower returns the borrower's ongoing (pending/approved/returned)
// request for a book, used to block duplicate applications and to render the borrower's view.
func (r *BorrowRepository) FindOngoingByBookAndBorrower(bookID, borrowerID uint) (*model.BorrowRequest, error) {
	var req model.BorrowRequest
	err := translate(r.db.
		Where("book_id = ? AND borrower_id = ? AND status IN ?", bookID, borrowerID,
			[]string{constants.BorrowStatusPending, constants.BorrowStatusApproved, constants.BorrowStatusReturned}).
		Order("id DESC").First(&req).Error)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// ListByBook returns all borrow requests of a book (seller review), pending first.
func (r *BorrowRepository) ListByBook(bookID uint) ([]model.BorrowRequest, error) {
	var items []model.BorrowRequest
	err := r.db.Preload("Borrower").
		Where("book_id = ?", bookID).
		Order("CASE WHEN status = 'pending' THEN 0 ELSE 1 END, created_at DESC").
		Find(&items).Error
	return items, err
}

// ListByUser returns requests where the user is borrower or seller, by role.
func (r *BorrowRepository) ListByUser(userID uint, role string, offset, limit int) ([]model.BorrowRequest, int64, error) {
	var items []model.BorrowRequest
	var total int64
	db := r.db.Model(&model.BorrowRequest{}).Preload("Book").Preload("Borrower").Preload("Seller")
	if role == "seller" {
		db = db.Where("seller_id = ?", userID)
	} else {
		db = db.Where("borrower_id = ?", userID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// RejectOtherPendingTx rejects every other pending request for the book
// (called when one request is approved).
func (r *BorrowRepository) RejectOtherPendingTx(tx *gorm.DB, bookID, exceptID uint) error {
	return translate(tx.Model(&model.BorrowRequest{}).
		Where("book_id = ? AND status = ? AND id <> ?", bookID, constants.BorrowStatusPending, exceptID).
		Update("status", constants.BorrowStatusRejected).Error)
}
