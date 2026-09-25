package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/model"
)

// MessageRepository handles message persistence.
type MessageRepository struct{ db *gorm.DB }

// NewMessageRepository creates a MessageRepository.
func NewMessageRepository(db *gorm.DB) *MessageRepository { return &MessageRepository{db: db} }

// Create inserts a message.
func (r *MessageRepository) Create(m *model.Message) error { return translate(r.db.Create(m).Error) }

// ListByConversation returns messages of a conversation ordered by time.
func (r *MessageRepository) ListByConversation(conversationID uint) ([]model.Message, error) {
	var items []model.Message
	if err := r.db.Preload("Sender").Where("conversation_id = ?", conversationID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CountUnread counts unread messages for a user in a conversation.
func (r *MessageRepository) CountUnread(conversationID, userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id <> ? AND is_read = ?", conversationID, userID, false).
		Count(&count).Error
	return count, err
}

// MarkConversationRead marks all messages as read for a user in a conversation.
func (r *MessageRepository) MarkConversationRead(conversationID, userID uint) error {
	return translate(r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id <> ? AND is_read = ?", conversationID, userID, false).
		Update("is_read", true).Error)
}
