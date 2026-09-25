package service

import (
	"fmt"
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

// BorrowService handles the short-term lending flow:
// apply -> approve/reject -> return -> confirm, plus overdue reminders.
type BorrowService struct {
	db          *gorm.DB
	borrowRepo  *repository.BorrowRepository
	bookRepo    *repository.BookRepository
	convService *ConversationService
	logger      *slog.Logger
}

// NewBorrowService creates a BorrowService.
func NewBorrowService(db *gorm.DB, borrowRepo *repository.BorrowRepository, bookRepo *repository.BookRepository,
	convService *ConversationService, logger *slog.Logger) *BorrowService {
	return &BorrowService{db: db, borrowRepo: borrowRepo, bookRepo: bookRepo, convService: convService, logger: logger}
}

// Apply creates a pending borrow request for a lendable, on-sale book.
func (s *BorrowService) Apply(borrowerID, bookID uint) (*dto.BorrowDTO, error) {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if book.SellerID == borrowerID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowSelf)
	}
	if !book.Lendable {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBookNotLendable, constants.MsgBookNotLendable)
	}
	if book.Status != constants.BookStatusOnSale {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid+": 书籍当前不可申请借阅")
	}
	if existing, err := s.borrowRepo.FindOngoingByBookAndBorrower(bookID, borrowerID); err == nil && existing != nil {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowDuplicate, constants.MsgBorrowDuplicate)
	}
	req := &model.BorrowRequest{
		BookID:     bookID,
		BorrowerID: borrowerID,
		SellerID:   book.SellerID,
		Status:     constants.BorrowStatusPending,
	}
	if err := s.borrowRepo.Create(req); err != nil {
		s.logger.Error("borrow request create failed", "book_id", bookID, "borrower", borrowerID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogBorrowApplied, "id", req.ID, "book_id", req.BookID, "borrower", req.BorrowerID)
	created, err := s.borrowRepo.FindByID(req.ID)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	d := dto.FromBorrow(created)
	return &d, nil
}

// Approve transitions pending -> approved, marks the book lent_out with a due
// date computed from the configured lend days, and rejects all other pending
// applications. Everything runs in one transaction with row locks.
func (s *BorrowService) Approve(sellerID, requestID uint) (*dto.BorrowDTO, error) {
	return s.transitionRequest(sellerID, requestID, requestSellerApproves)
}

// Reject transitions pending -> rejected by the seller.
func (s *BorrowService) Reject(sellerID, requestID uint) (*dto.BorrowDTO, error) {
	return s.transitionRequest(sellerID, requestID, requestSellerRejects)
}

// Return transitions approved -> returned (borrower hands the book back).
func (s *BorrowService) Return(borrowerID, requestID uint) (*dto.BorrowDTO, error) {
	return s.transitionRequest(borrowerID, requestID, requestBorrowerReturns)
}

// ConfirmReturn transitions returned -> completed (seller confirms the book is
// back) and restores the book to on_sale so it can be lent/sold again.
func (s *BorrowService) ConfirmReturn(sellerID, requestID uint) (*dto.BorrowDTO, error) {
	return s.transitionRequest(sellerID, requestID, requestSellerConfirms)
}

type requestAction int

const (
	requestSellerApproves requestAction = iota
	requestSellerRejects
	requestBorrowerReturns
	requestSellerConfirms
)

func (s *BorrowService) transitionRequest(userID, requestID uint, action requestAction) (*dto.BorrowDTO, error) {
	var result *dto.BorrowDTO
	err := s.db.Transaction(func(tx *gorm.DB) error {
		req, err := s.borrowRepo.FindByIDForUpdate(tx, requestID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgBorrowNotFound)
		}
		book, err := s.bookRepo.FindByIDForUpdate(tx, req.BookID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound)
		}
		now := time.Now()
		switch action {
		case requestSellerApproves:
			if req.SellerID != userID || book.SellerID != userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
			}
			if req.Status != constants.BorrowStatusPending {
				return util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid)
			}
			if book.Status != constants.BookStatusOnSale {
				return util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid+": 书籍当前状态不可借出")
			}
			lendDays := book.LendDays
			if lendDays != constants.LendDays7 && lendDays != constants.LendDays14 {
				lendDays = constants.LendDays7
			}
			due := now.Add(time.Duration(lendDays) * 24 * time.Hour)
			req.Status = constants.BorrowStatusApproved
			req.ApprovedAt = &now
			req.DueAt = &due
			book.Status = constants.BookStatusLentOut
		case requestSellerRejects:
			if req.SellerID != userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
			}
			if req.Status != constants.BorrowStatusPending {
				return util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid)
			}
			req.Status = constants.BorrowStatusRejected
		case requestBorrowerReturns:
			if req.BorrowerID != userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotBorrower)
			}
			if req.Status != constants.BorrowStatusApproved {
				return util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid)
			}
			req.Status = constants.BorrowStatusReturned
			req.ReturnedAt = &now
		case requestSellerConfirms:
			if req.SellerID != userID {
				return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
			}
			if req.Status != constants.BorrowStatusReturned {
				return util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid)
			}
			req.Status = constants.BorrowStatusCompleted
			req.ConfirmedAt = &now
			book.Status = constants.BookStatusOnSale
		}
		if err := s.borrowRepo.UpdateTx(tx, req); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		if err := s.bookRepo.UpdateTx(tx, book); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		if action == requestSellerApproves {
			if err := s.borrowRepo.RejectOtherPendingTx(tx, req.BookID, req.ID); err != nil {
				return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
			}
		}
		updated, err := s.borrowRepo.FindByIDTx(tx, req.ID)
		if err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		d := dto.FromBorrow(updated)
		result = &d
		return nil
	})
	if err != nil {
		s.logger.Warn("borrow request transition failed", "id", requestID, "action", action, "user_id", userID, "error", err)
		return nil, err
	}
	switch action {
	case requestSellerApproves:
		s.logger.Info(constants.LogBorrowApproved, "id", requestID, "book_id", result.BookID, "due_at", result.DueAt)
	case requestSellerRejects:
		s.logger.Info(constants.LogBorrowRejected, "id", requestID, "book_id", result.BookID, "seller", userID)
	case requestBorrowerReturns:
		s.logger.Info(constants.LogBorrowReturned, "id", requestID, "book_id", result.BookID, "borrower", userID)
	case requestSellerConfirms:
		s.logger.Info(constants.LogBorrowConfirmReturned, "id", requestID, "book_id", result.BookID, "seller", userID)
	}
	return result, nil
}

// Remind delivers a return reminder from the seller to the borrower through the
// in-app chat (conversation auto-created if needed).
func (s *BorrowService) Remind(sellerID, requestID uint) (*dto.BorrowDTO, error) {
	req, err := s.borrowRepo.FindByID(requestID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgBorrowNotFound)
	}
	if req.SellerID != sellerID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
	}
	if !req.IsLentOut() {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowStatusInvalid, constants.MsgBorrowStatusInvalid)
	}
	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound)
	}
	content := s.reminderContent(book.Title, req)
	if err := s.convService.SendBookMessage(req.BookID, sellerID, req.BorrowerID, content); err != nil {
		return nil, err
	}
	now := time.Now()
	req.RemindedAt = &now
	if err := s.borrowRepo.Update(req); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogBorrowReminded, "id", requestID, "book_id", req.BookID, "seller", sellerID)
	updated, err := s.borrowRepo.FindByID(requestID)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	d := dto.FromBorrow(updated)
	return &d, nil
}

// ListByBook returns the borrow requests of a book (seller only).
func (s *BorrowService) ListByBook(sellerID, bookID uint) ([]dto.BorrowDTO, error) {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if book.SellerID != sellerID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
	}
	items, err := s.borrowRepo.ListByBook(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BorrowDTO, 0, len(items))
	for i := range items {
		list = append(list, dto.FromBorrow(&items[i]))
	}
	return list, nil
}

// ListMine returns the caller's borrow requests as borrower or seller.
func (s *BorrowService) ListMine(userID uint, role string, page, size int) (*dto.PageData, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = util.DefaultPageSize
	}
	if size > util.MaxPageSize {
		size = util.MaxPageSize
	}
	items, total, err := s.borrowRepo.ListByUser(userID, role, (page-1)*size, size)
	if err != nil {
		s.logger.Error("borrow request list failed", "user_id", userID, "role", role, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BorrowDTO, 0, len(items))
	for i := range items {
		list = append(list, dto.FromBorrow(&items[i]))
	}
	return &dto.PageData{List: list, Total: total, Page: page, Size: size}, nil
}

func (s *BorrowService) reminderContent(bookTitle string, req *model.BorrowRequest) string {
	if req.DueAt == nil {
		return fmt.Sprintf("【还书提醒】你借阅的《%s》请尽快归还，谢谢！", bookTitle)
	}
	dueText := req.DueAt.In(time.Local).Format("01-02 15:04")
	if req.IsOverdue() {
		return fmt.Sprintf("【还书提醒】你借阅的《%s》已于 %s 到期，请尽快归还，谢谢！", bookTitle, dueText)
	}
	return fmt.Sprintf("【还书提醒】你借阅的《%s》将于 %s 到期，请按时归还，谢谢！", bookTitle, dueText)
}
