package model

import (
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
)

// User 学生/管理员账号。
type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StudentNo     string    `gorm:"size:32;uniqueIndex;not null" json:"student_no"`
	Email         string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"size:255;not null" json:"-"`
	Name          string    `gorm:"size:64" json:"name"`
	Department    string    `gorm:"size:64" json:"department"`
	Campus        string    `gorm:"size:64" json:"campus"`
	Contact       string    `gorm:"size:64" json:"contact"`
	AvatarURL     string    `gorm:"size:255" json:"avatar_url"`
	Role          string    `gorm:"size:16;default:student;index" json:"role"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	Status        string    `gorm:"size:16;default:active;index" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// IsAdmin reports whether the user has the admin role.
func (u *User) IsAdmin() bool { return u.Role == constants.RoleAdmin }
