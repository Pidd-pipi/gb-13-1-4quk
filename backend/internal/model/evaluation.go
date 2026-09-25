package model

import "time"

// Evaluation 交易完成后买卖双方互评记录。
type Evaluation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FromUserID uint      `gorm:"index:idx_eval_pair,unique;not null" json:"from_user_id"`
	ToUserID   uint      `gorm:"index:idx_eval_pair,unique;not null" json:"to_user_id"`
	BookID     uint      `gorm:"index:idx_eval_pair,unique;not null" json:"book_id"`
	Type       string    `gorm:"size:16;not null" json:"type"`
	Content    string    `gorm:"size:500" json:"content"`
	CreatedAt  time.Time `json:"created_at"`

	FromUser *User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   *User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
	Book     *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}
