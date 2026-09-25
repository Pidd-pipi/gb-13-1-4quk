package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// CreateConversationRequest 从书籍或求购发起会话。
type CreateConversationRequest struct {
	BookID   uint   `json:"book_id"`
	WishID   uint   `json:"wish_id"`
	ToUserID uint   `json:"to_user_id"` // 卖家主动联系某位同学（如借阅人）时使用
	Content  string `json:"content" binding:"required,max=2000"`
}

// SendMessageRequest 发送消息。
type SendMessageRequest struct {
	Content  string `json:"content" binding:"max=2000"`
	ImageURL string `json:"image_url" binding:"omitempty,url"`
}

// ConversationDTO 会话响应体。
type ConversationDTO struct {
	ID            uint     `json:"id"`
	BookID        uint     `json:"book_id"`
	WishID        uint     `json:"wish_id"`
	BuyerID       uint     `json:"buyer_id"`
	SellerID      uint     `json:"seller_id"`
	LastMessage   string   `json:"last_message"`
	LastMessageAt string   `json:"last_message_at"`
	CreatedAt     string   `json:"created_at"`
	Book          *BookDTO `json:"book,omitempty"`
	Buyer         *UserDTO `json:"buyer,omitempty"`
	Seller        *UserDTO `json:"seller,omitempty"`
	UnreadCount   int64    `json:"unread_count"`
}

// FromConversation converts a model.Conversation to ConversationDTO.
func FromConversation(c *model.Conversation) ConversationDTO {
	dto := ConversationDTO{
		ID:          c.ID,
		BookID:      c.BookID,
		WishID:      c.WishID,
		BuyerID:     c.BuyerID,
		SellerID:    c.SellerID,
		LastMessage: c.LastMessage,
		CreatedAt:   util.FormatTime(c.CreatedAt),
	}
	if c.LastMessageAt != nil {
		dto.LastMessageAt = util.FormatTime(*c.LastMessageAt)
	}
	if c.Book != nil {
		b := FromBook(c.Book)
		dto.Book = &b
	}
	if c.Buyer != nil {
		u := FromUser(c.Buyer)
		dto.Buyer = &u
	}
	if c.Seller != nil {
		u := FromUser(c.Seller)
		dto.Seller = &u
	}
	return dto
}

// MessageDTO 消息响应体。
type MessageDTO struct {
	ID             uint     `json:"id"`
	ConversationID uint     `json:"conversation_id"`
	SenderID       uint     `json:"sender_id"`
	Content        string   `json:"content"`
	ImageURL       string   `json:"image_url"`
	IsRead         bool     `json:"is_read"`
	CreatedAt      string   `json:"created_at"`
	Sender         *UserDTO `json:"sender,omitempty"`
}

// FromMessage converts a model.Message to MessageDTO.
func FromMessage(m *model.Message) MessageDTO {
	dto := MessageDTO{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Content:        m.Content,
		ImageURL:       m.ImageURL,
		IsRead:         m.IsRead,
		CreatedAt:      util.FormatTime(m.CreatedAt),
	}
	if m.Sender != nil {
		u := FromUser(m.Sender)
		dto.Sender = &u
	}
	return dto
}
