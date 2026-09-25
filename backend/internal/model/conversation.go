package model

import (
	"time"
)

// Conversation 买卖双方交易会话。
type Conversation struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	BookID        uint       `gorm:"index" json:"book_id"`
	WishID        uint       `gorm:"index" json:"wish_id"`
	BuyerID       uint       `gorm:"index;not null" json:"buyer_id"`
	SellerID      uint       `gorm:"index;not null" json:"seller_id"`
	LastMessage   string     `gorm:"size:2000" json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Book   *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Buyer  *User `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller *User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}
