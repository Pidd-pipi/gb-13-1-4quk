package model

import "time"

// Message 会话内的文字/图片消息。
type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID uint      `gorm:"index;not null" json:"conversation_id"`
	SenderID       uint      `gorm:"index;not null" json:"sender_id"`
	Content        string    `gorm:"size:2000" json:"content"`
	ImageURL       string    `gorm:"size:255" json:"image_url"`
	IsRead         bool      `gorm:"default:false" json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`

	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
