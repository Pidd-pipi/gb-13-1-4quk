package model

import "time"

// BrowseHistory 用户浏览历史。
type BrowseHistory struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"index;not null" json:"user_id"`
	BookID   uint      `gorm:"index;not null" json:"book_id"`
	ViewedAt time.Time `gorm:"index;not null" json:"viewed_at"`
}
