package model

import (
	"time"

	"gorm.io/datatypes"
)

// Book 闲置书籍。
type Book struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	SellerID        uint           `gorm:"index;not null" json:"seller_id"`
	Title           string         `gorm:"size:128;not null;index" json:"title"`
	Author          string         `gorm:"size:64" json:"author"`
	ISBN            string         `gorm:"size:32;index" json:"isbn"`
	CourseName      string         `gorm:"size:128" json:"course_name"`
	OriginalPrice   float64        `gorm:"type:decimal(10,2);default:0" json:"original_price"`
	Price           float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Condition       string         `gorm:"size:16;not null;index" json:"condition"`
	SubjectCategory string         `gorm:"size:16;not null;index" json:"subject_category"`
	TradeType       string         `gorm:"size:16;not null" json:"trade_type"`
	Campus          string         `gorm:"size:64" json:"campus"`
	Description     string         `gorm:"type:text" json:"description"`
	Images          datatypes.JSON `gorm:"type:json" json:"images"`
	Status          string         `gorm:"size:16;default:on_sale;index" json:"status"`
	ReservedBy      uint           `gorm:"default:0" json:"reserved_by"`
	ReservedAt      *time.Time     `json:"reserved_at"`
	// Borrowable 是否开启短借；BorrowDuration 短借时长（7/14 天），仅 Borrowable=true 时有效。
	Borrowable     bool      `gorm:"not null;default:false" json:"borrowable"`
	BorrowDuration int       `gorm:"not null;default:0" json:"borrow_duration"`
	ViewCount      int       `gorm:"default:0" json:"view_count"`
	FavoriteCount  int       `gorm:"default:0" json:"favorite_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Seller *User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}
