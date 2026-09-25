package dto

import (
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// CreateEvaluationRequest 创建交易评价。
type CreateEvaluationRequest struct {
	ToUserID uint   `json:"to_user_id" binding:"required"`
	BookID   uint   `json:"book_id" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=good neutral bad"`
	Content  string `json:"content" binding:"max=500"`
}

// EvaluationDTO 评价响应体。
type EvaluationDTO struct {
	ID         uint     `json:"id"`
	FromUserID uint     `json:"from_user_id"`
	ToUserID   uint     `json:"to_user_id"`
	BookID     uint     `json:"book_id"`
	Type       string   `json:"type"`
	TypeText   string   `json:"type_text"`
	Content    string   `json:"content"`
	CreatedAt  string   `json:"created_at"`
	FromUser   *UserDTO `json:"from_user,omitempty"`
	Book       *BookDTO `json:"book,omitempty"`
}

// FromEvaluation converts a model.Evaluation to EvaluationDTO.
func FromEvaluation(e *model.Evaluation) EvaluationDTO {
	dto := EvaluationDTO{
		ID:         e.ID,
		FromUserID: e.FromUserID,
		ToUserID:   e.ToUserID,
		BookID:     e.BookID,
		Type:       e.Type,
		TypeText:   util.FormatEvaluationTypeText(e.Type),
		Content:    e.Content,
		CreatedAt:  util.FormatTime(e.CreatedAt),
	}
	if e.FromUser != nil {
		u := FromUser(e.FromUser)
		dto.FromUser = &u
	}
	if e.Book != nil {
		b := FromBook(e.Book)
		dto.Book = &b
	}
	return dto
}
