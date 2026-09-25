package service

import (
	"encoding/json"

	"errors"
	"fmt"
	"gorm.io/datatypes"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/util"
)

// BookService handles book lifecycle: CRUD, search, favorites, history,
// recommendations and the on_sale -> reserved -> sold state machine.
type BookService struct {
	db          *gorm.DB
	bookRepo    *repository.BookRepository
	favRepo     *repository.FavoriteRepository
	historyRepo *repository.BrowseHistoryRepository
	userRepo    *repository.UserRepository
	logger      *slog.Logger
}

// NewBookService creates a BookService.
func NewBookService(db *gorm.DB, bookRepo *repository.BookRepository, favRepo *repository.FavoriteRepository,
	historyRepo *repository.BrowseHistoryRepository, userRepo *repository.UserRepository, logger *slog.Logger) *BookService {
	return &BookService{db: db, bookRepo: bookRepo, favRepo: favRepo, historyRepo: historyRepo, userRepo: userRepo, logger: logger}
}

// CreateBook publishes a new book for sale.
func (s *BookService) CreateBook(sellerID uint, req dto.CreateBookRequest) (*dto.BookDTO, error) {
	book := &model.Book{
		SellerID:        sellerID,
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		CourseName:      req.CourseName,
		OriginalPrice:   req.OriginalPrice,
		Price:           req.Price,
		Condition:       req.Condition,
		SubjectCategory: req.SubjectCategory,
		TradeType:       req.TradeType,
		Campus:          req.Campus,
		Description:     req.Description,
		Images:          marshalImages(req.Images),
		Status:          constants.BookStatusOnSale,
	}
	if err := s.bookRepo.Create(book); err != nil {
		s.logger.Error(constants.LogBookCreateFailed, "seller_id", sellerID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogBookCreateSuccess, "id", book.ID, "title", book.Title, "seller_id", sellerID)
	created, err := s.bookRepo.FindByID(book.ID)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	d := dto.FromBook(created)
	return &d, nil
}

// UpdateBook updates a book owned by the caller (only when on_sale).
func (s *BookService) UpdateBook(userID, bookID uint, req dto.UpdateBookRequest) (*dto.BookDTO, error) {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if book.SellerID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBookNotOwned, constants.MsgBookNotOwned)
	}
	if book.Status != constants.BookStatusOnSale {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBookStatusInvalid, constants.MsgBookStatusInvalid+": 仅"+"在售"+"书籍可编辑")
	}
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.ISBN != "" {
		book.ISBN = req.ISBN
	}
	if req.CourseName != "" {
		book.CourseName = req.CourseName
	}
	if req.Price > 0 {
		book.Price = req.Price
	}
	if req.Condition != "" {
		book.Condition = req.Condition
	}
	if req.SubjectCategory != "" {
		book.SubjectCategory = req.SubjectCategory
	}
	if req.TradeType != "" {
		book.TradeType = req.TradeType
	}
	if req.Campus != "" {
		book.Campus = req.Campus
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Images != nil {
		book.Images = marshalImages(req.Images)
	}
	if req.OriginalPrice >= 0 {
		book.OriginalPrice = req.OriginalPrice
	}
	if err := s.bookRepo.Update(book); err != nil {
		s.logger.Error(constants.LogBookUpdateSuccess, "id", bookID, "title", book.Title, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogBookUpdateSuccess, "id", bookID, "title", book.Title)
	d := dto.FromBook(book)
	return &d, nil
}

// DeleteBook removes a book owned by the caller.
func (s *BookService) DeleteBook(userID, bookID uint) error {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if book.SellerID != userID {
		return util.NewAppError(http.StatusForbidden, constants.CodeBookNotOwned, constants.MsgBookNotOwned)
	}
	if err := s.bookRepo.Delete(bookID); err != nil {
		s.logger.Error(constants.LogBookDeleteSuccess, "id", bookID, "seller_id", userID, "error", err)
		return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogBookDeleteSuccess, "id", bookID, "seller_id", userID)
	return nil
}

// ListBooks searches and filters books. Also reused by GetRecommendations
// for a narrower filter, keeping the filter pipeline in one place.
func (s *BookService) ListBooks(q dto.BookQuery) (*dto.PageData, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = util.DefaultPageSize
	}
	if q.PageSize > util.MaxPageSize {
		q.PageSize = util.MaxPageSize
	}
	items, total, err := s.bookRepo.List(q)
	if err != nil {
		s.logger.Error("book list failed", "keyword", q.Keyword, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BookDTO, 0, len(items))
	for _, b := range items {
		d := dto.FromBook(&b)
		if b.Seller != nil {
			u := dto.FromUser(b.Seller)
			d.Seller = &u
		}
		list = append(list, d)
	}
	return &dto.PageData{List: list, Total: total, Page: q.Page, Size: q.PageSize}, nil
}

// GetBookDetail returns a book, records view count and browse history.
func (s *BookService) GetBookDetail(userID, bookID uint) (*dto.BookDTO, error) {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if err := s.bookRepo.IncrementView(bookID); err != nil {
		s.logger.Warn("increment view failed", "id", bookID, "error", err)
	}
	book.ViewCount++
	if userID > 0 {
		_ = s.historyRepo.Create(&model.BrowseHistory{UserID: userID, BookID: bookID, ViewedAt: time.Now()})
	}
	s.logger.Info(constants.LogBookViewed, "id", bookID, "user_id", userID)
	d := dto.FromBook(book)
	if book.Seller != nil {
		u := dto.FromUser(book.Seller)
		d.Seller = &u
	}
	if userID > 0 {
		exists, _ := s.favRepo.Exists(userID, bookID)
		d.IsFavorite = exists
	}
	return &d, nil
}

// ReserveBook transitions on_sale -> reserved (buyer marks 已预约).
func (s *BookService) ReserveBook(userID, bookID uint) (*dto.BookDTO, error) {
	return s.transition(userID, bookID, constants.BookStatusOnSale, constants.BookStatusReserved)
}

// CancelReserve transitions reserved -> on_sale (seller or reserving buyer).
func (s *BookService) CancelReserve(userID, bookID uint) (*dto.BookDTO, error) {
	return s.transition(userID, bookID, constants.BookStatusReserved, constants.BookStatusOnSale)
}

// MarkSold transitions on_sale/reserved -> sold (seller confirms).
func (s *BookService) MarkSold(userID, bookID uint) (*dto.BookDTO, error) {
	book, err := s.transition(userID, bookID, "", constants.BookStatusSold)
	return book, err
}

// transition executes the book status machine inside a transaction using
// SELECT ... FOR UPDATE to prevent concurrent transitions.
func (s *BookService) transition(userID, bookID uint, from, to string) (*dto.BookDTO, error) {
	var result *dto.BookDTO
	err := s.db.Transaction(func(tx *gorm.DB) error {
		book, err := s.bookRepo.FindByIDForUpdate(tx, bookID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
		}
		fromStatus := book.Status
		if from != "" && fromStatus != from {
			return util.NewAppError(http.StatusConflict, constants.CodeBookStatusConflict,
				fmt.Sprintf("%s: 书籍(id=%d)当前状态=%s", constants.MsgBookStatusInvalid, bookID, util.FormatBookStatusText(fromStatus)))
		}
		switch to {
		case constants.BookStatusReserved:
			if fromStatus != constants.BookStatusOnSale {
				return util.NewAppError(http.StatusConflict, constants.CodeBookStatusConflict, constants.MsgBookStatusInvalid)
			}
			if book.SellerID == userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBookStatusInvalid, "不能预约自己发布的书籍")
			}
			now := time.Now()
			book.Status = constants.BookStatusReserved
			book.ReservedBy = userID
			book.ReservedAt = &now
		case constants.BookStatusOnSale:
			if book.SellerID != userID && book.ReservedBy != userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBookNotOwned, "仅卖家或预约买家可取消预约")
			}
			book.Status = constants.BookStatusOnSale
			book.ReservedBy = 0
			book.ReservedAt = nil
		case constants.BookStatusSold:
			if book.SellerID != userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBookNotOwned, constants.MsgBookNotOwned)
			}
			if fromStatus == constants.BookStatusSold {
				return util.NewAppError(http.StatusConflict, constants.CodeBookStatusConflict, "书籍已售出，状态不可再变更")
			}
			book.Status = constants.BookStatusSold
		default:
			return util.NewAppError(http.StatusBadRequest, constants.CodeBookStatusInvalid, constants.MsgBookStatusInvalid)
		}
		if err := s.bookRepo.UpdateTx(tx, book); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		updated, err := s.bookRepo.FindByIDTx(tx, book.ID)
		if err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		d := dto.FromBook(updated)
		if updated.Seller != nil {
			u := dto.FromUser(updated.Seller)
			d.Seller = &u
		}
		result = &d
		s.logger.Info(constants.LogBookStatusChanged, "id", bookID, "from", fromStatus, "to", to, "operator", userID)
		return nil
	})
	if err != nil {
		s.logger.Warn(constants.LogBookStatusChangeFail, "id", bookID, "to", to, "error", err)
		return nil, err
	}
	return result, nil
}

// AddFavorite favorites a book for the user.
func (s *BookService) AddFavorite(userID, bookID uint) error {
	if _, err := s.bookRepo.FindByID(bookID); err != nil {
		return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	fav := &model.Favorite{UserID: userID, BookID: bookID}
	if err := s.favRepo.Create(fav); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return util.NewAppError(http.StatusConflict, constants.CodeBookAlreadyFavorite, "该书籍已在收藏夹")
		}
		s.logger.Error(constants.LogBookFavoriteAdded, "user_id", userID, "book_id", bookID, "error", err)
		return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	_ = s.bookRepo.IncrementFavorite(bookID, 1)
	s.logger.Info(constants.LogBookFavoriteAdded, "user_id", userID, "book_id", bookID)
	return nil
}

// RemoveFavorite removes a favorite for the user.
func (s *BookService) RemoveFavorite(userID, bookID uint) error {
	if err := s.favRepo.DeleteByUserBook(userID, bookID); err != nil {
		return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFavorite, "该书籍不在收藏夹")
	}
	_ = s.bookRepo.IncrementFavorite(bookID, -1)
	s.logger.Info(constants.LogBookFavoriteRemoved, "user_id", userID, "book_id", bookID)
	return nil
}

// ListFavorites returns the user's favorited books.
func (s *BookService) ListFavorites(userID uint, page, size int) (*dto.PageData, error) {
	ids, total, err := s.favRepo.ListBookIDs(userID, (page-1)*size, size)
	if err != nil {
		s.logger.Error("favorite list failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	books, err := s.bookRepo.ListByIDs(ids)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BookDTO, 0, len(books))
	for _, b := range books {
		d := dto.FromBook(&b)
		d.IsFavorite = true
		if b.Seller != nil {
			u := dto.FromUser(b.Seller)
			d.Seller = &u
		}
		list = append(list, d)
	}
	return &dto.PageData{List: list, Total: total, Page: page, Size: size}, nil
}

// ListHistory returns the user's recently viewed books.
func (s *BookService) ListHistory(userID uint, limit int) (*dto.PageData, error) {
	if limit <= 0 {
		limit = 20
	}
	ids, err := s.historyRepo.ListBookIDs(userID, limit)
	if err != nil {
		s.logger.Error("history list failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	books, err := s.bookRepo.ListByIDs(ids)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BookDTO, 0, len(books))
	for _, b := range books {
		d := dto.FromBook(&b)
		if b.Seller != nil {
			u := dto.FromUser(b.Seller)
			d.Seller = &u
		}
		list = append(list, d)
	}
	return &dto.PageData{List: list, Total: int64(len(list)), Page: 1, Size: limit}, nil
}

// GetRecommendations returns books from users in the same department,
// reusing ListBooks-style DTO conversion (shared service method contract).
func (s *BookService) GetRecommendations(userID uint, limit int) ([]dto.BookDTO, error) {
	if limit <= 0 {
		limit = 10
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound+": user id="+fmt.Sprint(userID))
	}
	items, err := s.bookRepo.ListRecommendations(user.Department, userID, limit)
	if err != nil {
		s.logger.Error(constants.LogBookRecommendList, "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BookDTO, 0, len(items))
	for _, b := range items {
		d := dto.FromBook(&b)
		if b.Seller != nil {
			u := dto.FromUser(b.Seller)
			d.Seller = &u
		}
		list = append(list, d)
	}
	s.logger.Info(constants.LogBookRecommendList, "user_id", userID)
	return list, nil
}

func marshalImages(images []string) datatypes.JSON {
	if len(images) == 0 {
		return datatypes.JSON("[]")
	}
	b, _ := json.Marshal(images)
	return datatypes.JSON(b)
}
