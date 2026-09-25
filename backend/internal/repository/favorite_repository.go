package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/model"
)

// FavoriteRepository handles book favorites.
type FavoriteRepository struct{ db *gorm.DB }

// NewFavoriteRepository creates a FavoriteRepository.
func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository { return &FavoriteRepository{db: db} }

// Create inserts a favorite (unique user+book).
func (r *FavoriteRepository) Create(f *model.Favorite) error { return translate(r.db.Create(f).Error) }

// DeleteByUserBook removes a favorite, returning ErrNotFound when absent.
func (r *FavoriteRepository) DeleteByUserBook(userID, bookID uint) error {
	res := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&model.Favorite{})
	if res.Error != nil {
		return translate(res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Exists checks whether a favorite exists.
func (r *FavoriteRepository) Exists(userID, bookID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count).Error
	return count > 0, err
}

// ListBookIDs returns favorite book ids for a user ordered by time desc.
func (r *FavoriteRepository) ListBookIDs(userID uint, offset, limit int) ([]uint, int64, error) {
	var ids []uint
	var total int64
	q := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Favorite
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	for _, row := range rows {
		ids = append(ids, row.BookID)
	}
	return ids, total, nil
}
