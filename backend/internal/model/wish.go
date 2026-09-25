package model

import (
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
)

// Wish 求购信息。
type Wish struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	UserID               uint      `gorm:"index;not null" json:"user_id"`
	BookTitle            string    `gorm:"size:128;not null;index" json:"book_title"`
	Author               string    `gorm:"size:64" json:"author"`
	ISBN                 string    `gorm:"size:32" json:"isbn"`
	ExpectedPrice        float64   `gorm:"type:decimal(10,2);default:0" json:"expected_price"`
	ConditionRequirement string    `gorm:"size:16" json:"condition_requirement"`
	SubjectCategory      string    `gorm:"size:16;not null;index" json:"subject_category"`
	Description          string    `gorm:"type:text" json:"description"`
	Status               string    `gorm:"size:16;default:open;index" json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// IsOpen reports whether the wish is still open.
func (w *Wish) IsOpen() bool { return w.Status == constants.WishStatusOpen }
