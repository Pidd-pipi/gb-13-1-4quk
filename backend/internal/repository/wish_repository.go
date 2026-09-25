package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
)

// WishRepository handles wish persistence.
type WishRepository struct{ db *gorm.DB }

// NewWishRepository creates a WishRepository.
func NewWishRepository(db *gorm.DB) *WishRepository { return &WishRepository{db: db} }

// Create inserts a wish.
func (r *WishRepository) Create(w *model.Wish) error { return translate(r.db.Create(w).Error) }

// FindByID locates a wish by id with user preloaded.
func (r *WishRepository) FindByID(id uint) (*model.Wish, error) {
	var w model.Wish
	if err := translate(r.db.Preload("User").First(&w, id).Error); err != nil {
		return nil, err
	}
	return &w, nil
}

// Update persists a wish.
func (r *WishRepository) Update(w *model.Wish) error { return translate(r.db.Save(w).Error) }

// Delete removes a wish.
func (r *WishRepository) Delete(id uint) error {
	return translate(r.db.Delete(&model.Wish{}, id).Error)
}

// List returns wishes filtered by WishQuery.
func (r *WishRepository) List(q dto.WishQuery) ([]model.Wish, int64, error) {
	var items []model.Wish
	var total int64
	db := r.db.Model(&model.Wish{}).Preload("User")
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		db = db.Where("book_title LIKE ? OR author LIKE ? OR isbn LIKE ?", like, like, like)
	}
	if q.SubjectCategory != "" {
		db = db.Where("subject_category = ?", q.SubjectCategory)
	}
	if q.UserID > 0 {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset(q.PageSize * (q.Page - 1)).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
