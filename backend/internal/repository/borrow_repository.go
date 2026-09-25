package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
)

// BorrowRepository handles short-borrow request persistence.
type BorrowRepository struct{ db *gorm.DB }

// NewBorrowRepository creates a BorrowRepository.
func NewBorrowRepository(db *gorm.DB) *BorrowRepository { return &BorrowRepository{db: db} }

// Create inserts a borrow request.
func (r *BorrowRepository) Create(b *model.Borrow) error { return translate(r.db.Create(b).Error) }

// CreateTx inserts a borrow request inside a transaction.
func (r *BorrowRepository) CreateTx(tx *gorm.DB, b *model.Borrow) error {
	return translate(tx.Create(b).Error)
}

// FindByID locates a borrow request by id with relations preloaded.
func (r *BorrowRepository) FindByID(id uint) (*model.Borrow, error) {
	var b model.Borrow
	err := translate(r.db.Preload("Book").Preload("Book.Seller").
		Preload("Lender").Preload("Borrower").First(&b, id).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// FindByIDForUpdate locks the borrow row inside a state-transition transaction.
func (r *BorrowRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.Borrow, error) {
	var b model.Borrow
	if err := translate(tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, id).Error); err != nil {
		return nil, err
	}
	return &b, nil
}

// FindByIDTx locates a borrow request inside a transaction with relations preloaded.
func (r *BorrowRepository) FindByIDTx(tx *gorm.DB, id uint) (*model.Borrow, error) {
	var b model.Borrow
	err := translate(tx.Preload("Book").Preload("Book.Seller").
		Preload("Lender").Preload("Borrower").First(&b, id).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// UpdateTx persists a borrow request inside a transaction.
func (r *BorrowRepository) UpdateTx(tx *gorm.DB, b *model.Borrow) error {
	return translate(tx.Save(b).Error)
}

// Update persists a borrow request.
func (r *BorrowRepository) Update(b *model.Borrow) error { return translate(r.db.Save(b).Error) }

// List returns borrow requests filtered by BorrowQuery, newest first.
func (r *BorrowRepository) List(q dto.BorrowQuery) ([]model.Borrow, int64, error) {
	var items []model.Borrow
	var total int64
	db := r.db.Model(&model.Borrow{}).
		Preload("Book").Preload("Book.Seller").
		Preload("Lender").Preload("Borrower")
	if q.LenderID > 0 {
		db = db.Where("lender_id = ?", q.LenderID)
	}
	if q.BorrowerID > 0 {
		db = db.Where("borrower_id = ?", q.BorrowerID)
	}
	if q.BookID > 0 {
		db = db.Where("book_id = ?", q.BookID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, size := q.Page, q.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if err := db.Order("created_at DESC").Offset(size * (page - 1)).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// FindPendingByBorrower returns the borrower's still-pending request for a book.
func (r *BorrowRepository) FindPendingByBorrower(bookID, borrowerID uint) (*model.Borrow, error) {
	return r.FindPendingByBorrowerTx(r.db, bookID, borrowerID)
}

// FindPendingByBorrowerTx is the transaction-scoped variant.
func (r *BorrowRepository) FindPendingByBorrowerTx(tx *gorm.DB, bookID, borrowerID uint) (*model.Borrow, error) {
	var b model.Borrow
	err := translate(tx.Where("book_id = ? AND borrower_id = ? AND status = ?",
		bookID, borrowerID, constants.BorrowStatusPending).First(&b).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// FindActiveByBook returns the active (approved / returning) borrow of a book if any.
func (r *BorrowRepository) FindActiveByBook(tx *gorm.DB, bookID uint) (*model.Borrow, error) {
	var b model.Borrow
	err := tx.Where("book_id = ? AND status IN ?", bookID,
		[]string{constants.BorrowStatusApproved, constants.BorrowStatusReturning}).First(&b).Error
	if err != nil {
		return nil, translate(err)
	}
	return &b, nil
}

// FindActiveByBookFull returns the active borrow with relations preloaded.
func (r *BorrowRepository) FindActiveByBookFull(bookID uint) (*model.Borrow, error) {
	var b model.Borrow
	err := translate(r.db.Preload("Lender").Preload("Borrower").
		Where("book_id = ? AND status IN ?", bookID,
			[]string{constants.BorrowStatusApproved, constants.BorrowStatusReturning}).
		First(&b).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// RejectPendingByBookTx marks every other pending request for a book as rejected
// when a request is approved (a book can only be lent to one borrower).
func (r *BorrowRepository) RejectPendingByBookTx(tx *gorm.DB, bookID, approvedID uint, reason string) (int64, error) {
	res := tx.Model(&model.Borrow{}).
		Where("book_id = ? AND id <> ? AND status = ?", bookID, approvedID, constants.BorrowStatusPending).
		Updates(map[string]interface{}{"status": constants.BorrowStatusRejected, "reject_reason": reason})
	return res.RowsAffected, translate(res.Error)
}

// FindLatestByBorrower returns the borrower's most recent request for a book,
// used to render the current viewer's apply state on the book detail page.
func (r *BorrowRepository) FindLatestByBorrower(bookID, borrowerID uint) (*model.Borrow, error) {
	var b model.Borrow
	err := translate(r.db.Where("book_id = ? AND borrower_id = ?", bookID, borrowerID).
		Order("created_at DESC").First(&b).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// CountPendingByBook counts pending requests for a book (seller view).
func (r *BorrowRepository) CountPendingByBook(bookID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Borrow{}).
		Where("book_id = ? AND status = ?", bookID, constants.BorrowStatusPending).
		Count(&count).Error
	return count, err
}
