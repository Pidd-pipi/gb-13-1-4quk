package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// UpdateProfileRequest 完善/更新个人资料请求。
type UpdateProfileRequest struct {
	Name       string `json:"name" binding:"required,max=64"`
	Department string `json:"department" binding:"required,max=64"`
	Campus     string `json:"campus" binding:"required,max=64"`
	Contact    string `json:"contact" binding:"required,max=64"`
}

// UpdateAvatarRequest 更新头像请求。
type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" binding:"required,url"`
}

// UserDTO 用户响应体。
type UserDTO struct {
	ID            uint   `json:"id"`
	StudentNo     string `json:"student_no"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Department    string `json:"department"`
	Campus        string `json:"campus"`
	Contact       string `json:"contact"`
	AvatarURL     string `json:"avatar_url"`
	Role          string `json:"role"`
	EmailVerified bool   `json:"email_verified"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// FromUser converts a model.User to UserDTO.
func FromUser(u *model.User) UserDTO {
	return UserDTO{
		ID:            u.ID,
		StudentNo:     u.StudentNo,
		Email:         u.Email,
		Name:          u.Name,
		Department:    u.Department,
		Campus:        u.Campus,
		Contact:       u.Contact,
		AvatarURL:     u.AvatarURL,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		Status:        u.Status,
		CreatedAt:     util.FormatTime(u.CreatedAt),
	}
}

// UserStatsDTO 用户主页统计（好评率与风险提示）。
type UserStatsDTO struct {
	UserID           uint   `json:"user_id"`
	TotalBooks       int64  `json:"total_books"`
	OnSaleBooks      int64  `json:"on_sale_books"`
	SoldBooks        int64  `json:"sold_books"`
	TotalEvaluations int64  `json:"total_evaluations"`
	GoodCount        int64  `json:"good_count"`
	NeutralCount     int64  `json:"neutral_count"`
	BadCount         int64  `json:"bad_count"`
	GoodRate         string `json:"good_rate"`
	RiskFlagged      bool   `json:"risk_flagged"`
}
