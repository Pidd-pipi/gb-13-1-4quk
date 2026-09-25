package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/model"
)

// BrowseHistoryRepository handles browse history persistence.
type BrowseHistoryRepository struct{ db *gorm.DB }

// NewBrowseHistoryRepository creates a BrowseHistoryRepository.
func NewBrowseHistoryRepository(db *gorm.DB) *BrowseHistoryRepository {
	return &BrowseHistoryRepository{db: db}
}

// Create inserts a browse history record.
func (r *BrowseHistoryRepository) Create(h *model.BrowseHistory) error {
	return translate(r.db.Create(h).Error)
}

// ListBookIDs returns recently viewed book ids for a user.
func (r *BrowseHistoryRepository) ListBookIDs(userID uint, limit int) ([]uint, error) {
	var rows []model.BrowseHistory
	if err := r.db.Where("user_id = ?", userID).Order("viewed_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.BookID)
	}
	return ids, nil
}
