package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
)

// BookRepository handles book persistence.
type BookRepository struct{ db *gorm.DB }

// NewBookRepository creates a BookRepository.
func NewBookRepository(db *gorm.DB) *BookRepository { return &BookRepository{db: db} }

// Create inserts a book.
func (r *BookRepository) Create(b *model.Book) error { return translate(r.db.Create(b).Error) }

// FindByID locates a book by id with seller preloaded.
func (r *BookRepository) FindByID(id uint) (*model.Book, error) {
	var b model.Book
	err := translate(r.db.Preload("Seller").First(&b, id).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// FindByIDForUpdate locks a book row for status transitions.
func (r *BookRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.Book, error) {
	var b model.Book
	err := translate(tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&b, id).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Update persists a book.
func (r *BookRepository) Update(b *model.Book) error { return translate(r.db.Save(b).Error) }

// UpdateTx persists a book inside a transaction.
func (r *BookRepository) UpdateTx(tx *gorm.DB, b *model.Book) error {
	return translate(tx.Save(b).Error)
}

// FindByIDTx returns a book inside a transaction with seller preloaded.
func (r *BookRepository) FindByIDTx(tx *gorm.DB, id uint) (*model.Book, error) {
	var b model.Book
	err := translate(tx.Preload("Seller").First(&b, id).Error)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Delete removes a book.
func (r *BookRepository) Delete(id uint) error {
	return translate(r.db.Delete(&model.Book{}, id).Error)
}

// List returns books filtered by BookQuery.
func (r *BookRepository) List(q dto.BookQuery) ([]model.Book, int64, error) {
	var items []model.Book
	var total int64
	db := r.db.Model(&model.Book{}).Preload("Seller")
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		db = db.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ? OR course_name LIKE ?", like, like, like, like)
	}
	if q.SubjectCategory != "" {
		db = db.Where("subject_category = ?", q.SubjectCategory)
	}
	if q.Condition != "" {
		db = db.Where("`condition` = ?", q.Condition)
	}
	if q.TradeType != "" {
		db = db.Where("trade_type = ?", q.TradeType)
	}
	if q.MinPrice > 0 {
		db = db.Where("price >= ?", q.MinPrice)
	}
	if q.MaxPrice > 0 {
		db = db.Where("price <= ?", q.MaxPrice)
	}
	if q.SellerID > 0 {
		db = db.Where("seller_id = ?", q.SellerID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.Borrowable != nil {
		db = db.Where("borrowable = ?", *q.Borrowable)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := "created_at DESC"
	switch q.Sort {
	case "price_asc":
		order = "price ASC"
	case "price_desc":
		order = "price DESC"
	case "newest":
		order = "created_at DESC"
	case "most_viewed":
		order = "view_count DESC"
	}
	if err := db.Order(order).Offset(q.PageSize * (q.Page - 1)).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListRecommendations returns books sold by users in the same department.
func (r *BookRepository) ListRecommendations(department string, excludeUser uint, limit int) ([]model.Book, error) {
	var items []model.Book
	err := r.db.Model(&model.Book{}).
		Preload("Seller").
		Joins("JOIN users ON users.id = books.seller_id").
		Where("users.department = ? AND books.seller_id <> ? AND books.status = ?", department, excludeUser, "on_sale").
		Order("books.created_at DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// IncrementView increments view count.
func (r *BookRepository) IncrementView(id uint) error {
	return translate(r.db.Model(&model.Book{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error)
}

// IncrementFavorite adjusts favorite count by delta.
func (r *BookRepository) IncrementFavorite(id uint, delta int) error {
	return translate(r.db.Model(&model.Book{}).Where("id = ?", id).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", delta)).Error)
}

// ListByIDs returns books for given ids (history/favorites hydration).
func (r *BookRepository) ListByIDs(ids []uint) ([]model.Book, error) {
	if len(ids) == 0 {
		return []model.Book{}, nil
	}
	var items []model.Book
	if err := r.db.Preload("Seller").Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
