package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// CreateBookRequest 发布书籍请求。
type CreateBookRequest struct {
	Title           string   `json:"title" binding:"required,max=128"`
	Author          string   `json:"author" binding:"max=64"`
	ISBN            string   `json:"isbn" binding:"max=32"`
	CourseName      string   `json:"course_name" binding:"max=128"`
	OriginalPrice   float64  `json:"original_price" binding:"gte=0"`
	Price           float64  `json:"price" binding:"required,gt=0"`
	Condition       string   `json:"condition" binding:"required,oneof=brand_new nine_new seven_new five_new"`
	SubjectCategory string   `json:"subject_category" binding:"required,oneof=science humanities econ_management art other"`
	TradeType       string   `json:"trade_type" binding:"required,oneof=in_person mail"`
	Campus          string   `json:"campus" binding:"max=64"`
	Description     string   `json:"description" binding:"max=2000"`
	Images          []string `json:"images" binding:"max=5"`
	Lendable        bool     `json:"lendable"`
	LendDays        int      `json:"lend_days" binding:"omitempty,oneof=7 14"`
}

// UpdateBookRequest 更新书籍请求。
type UpdateBookRequest struct {
	Title           string   `json:"title" binding:"max=128"`
	Author          string   `json:"author" binding:"max=64"`
	ISBN            string   `json:"isbn" binding:"max=32"`
	CourseName      string   `json:"course_name" binding:"max=128"`
	OriginalPrice   float64  `json:"original_price" binding:"gte=0"`
	Price           float64  `json:"price" binding:"gt=0"`
	Condition       string   `json:"condition" binding:"oneof=brand_new nine_new seven_new five_new"`
	SubjectCategory string   `json:"subject_category" binding:"oneof=science humanities econ_management art other"`
	TradeType       string   `json:"trade_type" binding:"oneof=in_person mail"`
	Campus          string   `json:"campus" binding:"max=64"`
	Description     string   `json:"description" binding:"max=2000"`
	Images          []string `json:"images" binding:"max=5"`
	Lendable        *bool    `json:"lendable"`
	LendDays        int      `json:"lend_days" binding:"omitempty,oneof=7 14"`
}

// BookQuery 书籍列表查询参数。
type BookQuery struct {
	Keyword         string  `form:"keyword"`
	SubjectCategory string  `form:"subject_category"`
	Condition       string  `form:"condition"`
	TradeType       string  `form:"trade_type"`
	MinPrice        float64 `form:"min_price"`
	MaxPrice        float64 `form:"max_price"`
	SellerID        uint    `form:"seller_id"`
	Status          string  `form:"status"` // on_sale / reserved / sold
	Sort            string  `form:"sort"`   // price_asc / price_desc / newest / most_viewed
	Page            int     `form:"page"`
	PageSize        int     `form:"page_size"`
}

// BookDTO 书籍响应体。
type BookDTO struct {
	ID              uint       `json:"id"`
	SellerID        uint       `json:"seller_id"`
	Title           string     `json:"title"`
	Author          string     `json:"author"`
	ISBN            string     `json:"isbn"`
	CourseName      string     `json:"course_name"`
	OriginalPrice   float64    `json:"original_price"`
	Price           float64    `json:"price"`
	Condition       string     `json:"condition"`
	ConditionText   string     `json:"condition_text"`
	SubjectCategory string     `json:"subject_category"`
	SubjectText     string     `json:"subject_text"`
	TradeType       string     `json:"trade_type"`
	TradeTypeText   string     `json:"trade_type_text"`
	Campus          string     `json:"campus"`
	Description     string     `json:"description"`
	Images          []string   `json:"images"`
	Status          string     `json:"status"`
	StatusText      string     `json:"status_text"`
	ReservedBy      uint       `json:"reserved_by"`
	Lendable        bool       `json:"lendable"`
	LendDays        int        `json:"lend_days"`
	ActiveBorrow    *BorrowDTO `json:"active_borrow,omitempty"`
	MyBorrow        *BorrowDTO `json:"my_borrow,omitempty"`
	ViewCount       int        `json:"view_count"`
	FavoriteCount   int        `json:"favorite_count"`
	CreatedAt       string     `json:"created_at"`
	Seller          *UserDTO   `json:"seller,omitempty"`
	IsFavorite      bool       `json:"is_favorite"`
}

// FromBook converts a model.Book to BookDTO.
func FromBook(b *model.Book) BookDTO {
	return BookDTO{
		ID:              b.ID,
		SellerID:        b.SellerID,
		Title:           b.Title,
		Author:          b.Author,
		ISBN:            b.ISBN,
		CourseName:      b.CourseName,
		OriginalPrice:   b.OriginalPrice,
		Price:           b.Price,
		Condition:       b.Condition,
		ConditionText:   util.FormatConditionText(b.Condition),
		SubjectCategory: b.SubjectCategory,
		SubjectText:     util.FormatSubjectText(b.SubjectCategory),
		TradeType:       b.TradeType,
		TradeTypeText:   util.FormatTradeTypeText(b.TradeType),
		Campus:          b.Campus,
		Description:     b.Description,
		Images:          util.SplitImages(string(b.Images)),
		Status:          b.Status,
		StatusText:      util.FormatBookStatusText(b.Status),
		ReservedBy:      b.ReservedBy,
		Lendable:        b.Lendable,
		LendDays:        b.LendDays,
		ViewCount:       b.ViewCount,
		FavoriteCount:   b.FavoriteCount,
		CreatedAt:       util.FormatTime(b.CreatedAt),
	}
}

// BookListDTO 书籍列表项。
type BookListDTO struct {
	BookDTO
}
