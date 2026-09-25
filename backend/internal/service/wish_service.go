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

// WishService handles wish (求购) lifecycle.
type WishService struct {
	wishRepo *repository.WishRepository
	logger   *slog.Logger
}

// NewWishService creates a WishService.
func NewWishService(wishRepo *repository.WishRepository, logger *slog.Logger) *WishService {
	return &WishService{wishRepo: wishRepo, logger: logger}
}

// CreateWish publishes a wish.
func (s *WishService) CreateWish(userID uint, req dto.CreateWishRequest) (*dto.WishDTO, error) {
	wish := &model.Wish{
		UserID:               userID,
		BookTitle:            req.BookTitle,
		Author:               req.Author,
		ISBN:                 req.ISBN,
		ExpectedPrice:        req.ExpectedPrice,
		ConditionRequirement: req.ConditionRequirement,
		SubjectCategory:      req.SubjectCategory,
		Description:          req.Description,
		Status:               constants.WishStatusOpen,
	}
	if err := s.wishRepo.Create(wish); err != nil {
		s.logger.Error(constants.LogWishCreateSuccess, "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogWishCreateSuccess, "id", wish.ID, "book_title", wish.BookTitle, "user_id", userID)
	created, err := s.wishRepo.FindByID(wish.ID)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	d := dto.FromWish(created)
	return &d, nil
}

// UpdateWish updates an owned open wish.
func (s *WishService) UpdateWish(userID, wishID uint, req dto.UpdateWishRequest) (*dto.WishDTO, error) {
	wish, err := s.wishRepo.FindByID(wishID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeWishNotFound, constants.MsgNotFound+": wish id="+fmt.Sprint(wishID))
	}
	if wish.UserID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeWishNotOwned, constants.MsgWishNotOwned)
	}
	if !wish.IsOpen() {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeWishAlreadyClosed, constants.MsgWishNotOwned)
	}
	if req.BookTitle != "" {
		wish.BookTitle = req.BookTitle
	}
	if req.Author != "" {
		wish.Author = req.Author
	}
	if req.ISBN != "" {
		wish.ISBN = req.ISBN
	}
	if req.ExpectedPrice >= 0 {
		wish.ExpectedPrice = req.ExpectedPrice
	}
	if req.ConditionRequirement != "" {
		wish.ConditionRequirement = req.ConditionRequirement
	}
	if req.SubjectCategory != "" {
		wish.SubjectCategory = req.SubjectCategory
	}
	if req.Description != "" {
		wish.Description = req.Description
	}
	if err := s.wishRepo.Update(wish); err != nil {
		s.logger.Error("wish update failed", "wish_id", wishID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	d := dto.FromWish(wish)
	return &d, nil
}

// CloseWish marks a wish as closed (owner).
func (s *WishService) CloseWish(userID, wishID uint) (*dto.WishDTO, error) {
	wish, err := s.wishRepo.FindByID(wishID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeWishNotFound, constants.MsgNotFound+": wish id="+fmt.Sprint(wishID))
	}
	if wish.UserID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeWishNotOwned, constants.MsgWishNotOwned)
	}
	if !wish.IsOpen() {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeWishAlreadyClosed, "求购信息已关闭")
	}
	wish.Status = constants.WishStatusClosed
	if err := s.wishRepo.Update(wish); err != nil {
		s.logger.Error(constants.LogWishCloseSuccess, "id", wishID, "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogWishCloseSuccess, "id", wishID, "user_id", userID)
	d := dto.FromWish(wish)
	return &d, nil
}

// DeleteWish removes an owned wish.
func (s *WishService) DeleteWish(userID, wishID uint) error {
	wish, err := s.wishRepo.FindByID(wishID)
	if err != nil {
		return util.NewAppError(http.StatusNotFound, constants.CodeWishNotFound, constants.MsgNotFound+": wish id="+fmt.Sprint(wishID))
	}
	if wish.UserID != userID {
		return util.NewAppError(http.StatusForbidden, constants.CodeWishNotOwned, constants.MsgWishNotOwned)
	}
	if err := s.wishRepo.Delete(wishID); err != nil {
		s.logger.Error("wish delete failed", "wish_id", wishID, "error", err)
		return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	return nil
}

// ListWishes searches wishes by keyword/category/user/status.
func (s *WishService) ListWishes(q dto.WishQuery) (*dto.PageData, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = util.DefaultPageSize
	}
	if q.PageSize > util.MaxPageSize {
		q.PageSize = util.MaxPageSize
	}
	items, total, err := s.wishRepo.List(q)
	if err != nil {
		s.logger.Error("wish list failed", "keyword", q.Keyword, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.WishDTO, 0, len(items))
	for _, w := range items {
		list = append(list, dto.FromWish(&w))
	}
	return &dto.PageData{List: list, Total: total, Page: q.Page, Size: q.PageSize}, nil
}

// GetWish returns a single wish.
func (s *WishService) GetWish(wishID uint) (*dto.WishDTO, error) {
	wish, err := s.wishRepo.FindByID(wishID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeWishNotFound, constants.MsgNotFound+": wish id="+fmt.Sprint(wishID))
	}
	d := dto.FromWish(wish)
	return &d, nil
}
