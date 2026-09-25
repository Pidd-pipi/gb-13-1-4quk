package model

import "time"

// Favorite 书籍收藏。
type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index:idx_fav_user_book,unique;not null" json:"user_id"`
	BookID    uint      `gorm:"index:idx_fav_user_book,unique;not null" json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}
