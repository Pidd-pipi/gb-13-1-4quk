package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/model"
)

// ConversationRepository handles conversation persistence.
type ConversationRepository struct{ db *gorm.DB }

// NewConversationRepository creates a ConversationRepository.
func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// Create inserts a conversation.
func (r *ConversationRepository) Create(c *model.Conversation) error {
	return translate(r.db.Create(c).Error)
}

// FindByID locates a conversation by id with relations preloaded.
func (r *ConversationRepository) FindByID(id uint) (*model.Conversation, error) {
	var c model.Conversation
	err := translate(r.db.Preload("Book").Preload("Buyer").Preload("Seller").First(&c, id).Error)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// FindExisting locates an existing conversation between two users for a book/wish.
func (r *ConversationRepository) FindExisting(bookID, wishID, buyerID, sellerID uint) (*model.Conversation, error) {
	var c model.Conversation
	q := r.db.Where("buyer_id = ? AND seller_id = ?", buyerID, sellerID)
	if bookID > 0 {
		q = q.Where("book_id = ?", bookID)
	}
	if wishID > 0 {
		q = q.Where("wish_id = ?", wishID)
	}
	err := translate(q.First(&c).Error)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListByUser returns conversations where the user is a participant.
func (r *ConversationRepository) ListByUser(userID uint, offset, limit int) ([]model.Conversation, int64, error) {
	var items []model.Conversation
	var total int64
	q := r.db.Model(&model.Conversation{}).
		Preload("Book").
		Preload("Buyer").
		Preload("Seller").
		Where("buyer_id = ? OR seller_id = ?", userID, userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Touch updates the last message preview and updated_at.
func (r *ConversationRepository) Touch(id uint, lastMessage string) error {
	return translate(r.db.Model(&model.Conversation{}).Where("id = ?", id).
		Updates(map[string]interface{}{"last_message": lastMessage, "last_message_at": time.Now()}).Error)
}
