package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/model"
)

// EvaluationRepository handles evaluation persistence.
type EvaluationRepository struct{ db *gorm.DB }

// NewEvaluationRepository creates an EvaluationRepository.
func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{db: db}
}

// Create inserts an evaluation.
func (r *EvaluationRepository) Create(e *model.Evaluation) error {
	return translate(r.db.Create(e).Error)
}

// ListByUser returns evaluations received by a user (paginated).
func (r *EvaluationRepository) ListByUser(userID uint, offset, limit int) ([]model.Evaluation, int64, error) {
	var items []model.Evaluation
	var total int64
	q := r.db.Model(&model.Evaluation{}).
		Preload("FromUser").
		Preload("Book").
		Where("to_user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// CountByType counts evaluations of a type received by a user.
func (r *EvaluationRepository) CountByType(userID, evalType string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Evaluation{}).
		Where("to_user_id = ? AND type = ?", userID, evalType).Count(&count).Error
	return count, err
}

// ExistsPair reports whether an evaluation from->to for a book already exists.
func (r *EvaluationRepository) ExistsPair(fromUserID, toUserID, bookID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Evaluation{}).
		Where("from_user_id = ? AND to_user_id = ? AND book_id = ?", fromUserID, toUserID, bookID).
		Count(&count).Error
	return count > 0, err
}
