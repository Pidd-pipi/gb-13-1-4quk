package service

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/util"
)

// EvaluationService handles post-trade evaluations.
type EvaluationService struct {
	evalRepo *repository.EvaluationRepository
	bookRepo *repository.BookRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewEvaluationService creates an EvaluationService.
func NewEvaluationService(evalRepo *repository.EvaluationRepository, bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *EvaluationService {
	return &EvaluationService{evalRepo: evalRepo, bookRepo: bookRepo, userRepo: userRepo, logger: logger}
}

// CreateEvaluation lets a buyer/seller evaluate the counterpart after a sold book.
func (s *EvaluationService) CreateEvaluation(fromUserID uint, req dto.CreateEvaluationRequest) (*dto.EvaluationDTO, error) {
	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(req.BookID))
	}
	if book.Status != constants.BookStatusSold {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeEvaluationInvalid, constants.MsgEvaluationForbidden)
	}
	isSeller := book.SellerID == fromUserID
	isBuyer := book.ReservedBy == fromUserID
	if !isSeller && !isBuyer {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeEvaluationForbidden, "仅交易双方可以评价")
	}
	if req.ToUserID != book.SellerID && req.ToUserID != book.ReservedBy {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeEvaluationForbidden, "评价对象不是本笔交易的参与者")
	}
	if req.ToUserID == fromUserID {
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeEvaluationInvalid, "不能评价自己")
	}
	exists, err := s.evalRepo.ExistsPair(fromUserID, req.ToUserID, req.BookID)
	if err != nil {
		s.logger.Error("evaluation exists check failed", "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	if exists {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeEvaluationDuplicate, constants.MsgEvaluationDuplicate)
	}
	eval := &model.Evaluation{
		FromUserID: fromUserID,
		ToUserID:   req.ToUserID,
		BookID:     req.BookID,
		Type:       req.Type,
		Content:    req.Content,
	}
	if err := s.evalRepo.Create(eval); err != nil {
		s.logger.Error(constants.LogEvaluationCreate, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogEvaluationCreate, "from", fromUserID, "to", req.ToUserID, "book_id", req.BookID, "type", req.Type)
	d := dto.FromEvaluation(eval)
	return &d, nil
}

// ListEvaluations returns evaluations received by a user. Also reused by the
// user stats page to hydrate the evaluation history list.
func (s *EvaluationService) ListEvaluations(userID uint, page, size int) (*dto.PageData, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = util.DefaultPageSize
	}
	items, total, err := s.evalRepo.ListByUser(userID, (page-1)*size, size)
	if err != nil {
		s.logger.Error("evaluation list failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.EvaluationDTO, 0, len(items))
	for i := range items {
		list = append(list, dto.FromEvaluation(&items[i]))
	}
	return &dto.PageData{List: list, Total: total, Page: page, Size: size}, nil
}
