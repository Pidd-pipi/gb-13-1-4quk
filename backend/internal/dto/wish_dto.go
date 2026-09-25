package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// CreateWishRequest 发布求购请求。
type CreateWishRequest struct {
	BookTitle            string  `json:"book_title" binding:"required,max=128"`
	Author               string  `json:"author" binding:"max=64"`
	ISBN                 string  `json:"isbn" binding:"max=32"`
	ExpectedPrice        float64 `json:"expected_price" binding:"gte=0"`
	ConditionRequirement string  `json:"condition_requirement" binding:"oneof=brand_new nine_new seven_new five_new"`
	SubjectCategory      string  `json:"subject_category" binding:"required,oneof=science humanities econ_management art other"`
	Description          string  `json:"description" binding:"max=1000"`
}

// UpdateWishRequest 更新求购请求。
type UpdateWishRequest struct {
	BookTitle            string  `json:"book_title" binding:"max=128"`
	Author               string  `json:"author" binding:"max=64"`
	ISBN                 string  `json:"isbn" binding:"max=32"`
	ExpectedPrice        float64 `json:"expected_price" binding:"gte=0"`
	ConditionRequirement string  `json:"condition_requirement" binding:"oneof=brand_new nine_new seven_new five_new"`
	SubjectCategory      string  `json:"subject_category" binding:"oneof=science humanities econ_management art other"`
	Description          string  `json:"description" binding:"max=1000"`
}

// WishQuery 求购列表查询参数。
type WishQuery struct {
	Keyword         string `form:"keyword"`
	SubjectCategory string `form:"subject_category"`
	UserID          uint   `form:"user_id"`
	Status          string `form:"status"`
	Page            int    `form:"page"`
	PageSize        int    `form:"page_size"`
}

// WishDTO 求购信息响应体。
type WishDTO struct {
	ID                   uint     `json:"id"`
	UserID               uint     `json:"user_id"`
	BookTitle            string   `json:"book_title"`
	Author               string   `json:"author"`
	ISBN                 string   `json:"isbn"`
	ExpectedPrice        float64  `json:"expected_price"`
	ConditionRequirement string   `json:"condition_requirement"`
	ConditionText        string   `json:"condition_text"`
	SubjectCategory      string   `json:"subject_category"`
	SubjectText          string   `json:"subject_text"`
	Description          string   `json:"description"`
	Status               string   `json:"status"`
	StatusText           string   `json:"status_text"`
	CreatedAt            string   `json:"created_at"`
	User                 *UserDTO `json:"user,omitempty"`
}

// FromWish converts a model.Wish to WishDTO.
func FromWish(w *model.Wish) WishDTO {
	dto := WishDTO{
		ID:                   w.ID,
		UserID:               w.UserID,
		BookTitle:            w.BookTitle,
		Author:               w.Author,
		ISBN:                 w.ISBN,
		ExpectedPrice:        w.ExpectedPrice,
		ConditionRequirement: w.ConditionRequirement,
		ConditionText:        util.FormatConditionText(w.ConditionRequirement),
		SubjectCategory:      w.SubjectCategory,
		SubjectText:          util.FormatSubjectText(w.SubjectCategory),
		Description:          w.Description,
		Status:               w.Status,
		StatusText:           util.FormatWishStatusText(w.Status),
		CreatedAt:            util.FormatTime(w.CreatedAt),
	}
	if w.User != nil {
		u := FromUser(w.User)
		dto.User = &u
	}
	return dto
}
