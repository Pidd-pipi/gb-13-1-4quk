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

// reminderInterval limits how often a seller can send an overdue reminder.
const reminderInterval = time.Minute

// defaultReminderContent is the reminder message body when the seller does not
// customize one. Reuses the in-app conversation/message channel.
func defaultReminderContent(bookTitle string, dueAt time.Time) string {
	return fmt.Sprintf("你好，你借阅的《%s》已于 %s 到期，请尽快归还，谢谢！", bookTitle, util.FormatTime(dueAt))
}

// BorrowService implements the short-borrow flow:
// apply(pending) -> approve(approved, book loaned, others rejected) ->
// return(returning) -> confirm(returned, book on_sale again); reject and
// overdue reminder are seller actions.
type BorrowService struct {
	db         *gorm.DB
	borrowRepo *repository.BorrowRepository
	bookRepo   *repository.BookRepository
	convRepo   *repository.ConversationRepository
	msgRepo    *repository.MessageRepository
	logger     *slog.Logger
}

// NewBorrowService creates a BorrowService.
func NewBorrowService(db *gorm.DB, borrowRepo *repository.BorrowRepository, bookRepo *repository.BookRepository,
	convRepo *repository.ConversationRepository, msgRepo *repository.MessageRepository, logger *slog.Logger) *BorrowService {
	return &BorrowService{db: db, borrowRepo: borrowRepo, bookRepo: bookRepo,
		convRepo: convRepo, msgRepo: msgRepo, logger: logger}
}

// Apply creates a pending borrow request from the book detail page.
func (s *BorrowService) Apply(borrowerID uint, bookID uint, req dto.CreateBorrowRequest) (*dto.BorrowDTO, error) {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if !book.Borrowable {
		return nil, util.NewAppError(http.StatusUnprocessableEntity, constants.CodeBorrowNotAllowed, constants.MsgBorrowNotAllowed)
	}
	if book.SellerID == borrowerID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, "不能借阅自己发布的书籍")
	}
	if book.Status != constants.BookStatusOnSale {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, "书籍当前不可借阅："+util.FormatBookStatusText(book.Status))
	}

	var result *dto.BorrowDTO
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 一本书同时只能有一笔生效借阅
		if active, e := s.borrowRepo.FindActiveByBook(tx, bookID); e == nil && active != nil {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, "该书已被借出，请稍后再试")
		}
		// 同一借阅人已有待处理申请
		if existing, e := s.borrowRepo.FindPendingByBorrowerTx(tx, bookID, borrowerID); e == nil && existing != nil {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowDuplicate, constants.MsgBorrowDuplicate)
		}
		b := &model.Borrow{
			BookID:       bookID,
			LenderID:     book.SellerID,
			BorrowerID:   borrowerID,
			Duration:     book.BorrowDuration,
			Status:       constants.BorrowStatusPending,
			RejectReason: "",
		}
		if err := s.borrowRepo.CreateTx(tx, b); err != nil {
			s.logger.Error(constants.LogBorrowApplied, "error", err)
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		full, err := s.borrowRepo.FindByIDTx(tx, b.ID)
		if err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		d := s.toDTO(full)
		result = &d
		s.logger.Info(constants.LogBorrowApplied, "id", b.ID, "book_id", bookID, "borrower", borrowerID, "duration", b.Duration)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Approve lets the seller accept a pending request: due date is recorded, the
// book becomes loaned, and all other pending requests are auto-rejected.
func (s *BorrowService) Approve(lenderID, borrowID uint) (*dto.BorrowDTO, error) {
	var result *dto.BorrowDTO
	err := s.db.Transaction(func(tx *gorm.DB) error {
		b, err := s.borrowRepo.FindByIDForUpdate(tx, borrowID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
		}
		if b.LenderID != lenderID {
			return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
		}
		if b.Status != constants.BorrowStatusPending {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, constants.MsgBorrowStatusInvalid)
		}
		book, err := s.bookRepo.FindByIDForUpdate(tx, b.BookID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound)
		}
		if active, e := s.borrowRepo.FindActiveByBook(tx, book.ID); e == nil && active != nil && active.ID != borrowID {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, "该书已有进行中的借阅")
		}
		if book.Status != constants.BookStatusOnSale {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, "书籍当前为"+util.FormatBookStatusText(book.Status)+"，无法同意借阅")
		}

		now := time.Now()
		due := now.AddDate(0, 0, b.Duration)
		b.Status = constants.BorrowStatusApproved
		b.ApprovedAt = &now
		b.DueAt = &due
		if err := s.borrowRepo.UpdateTx(tx, b); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}

		// 同一本书的其他待处理申请自动拒绝
		if _, err := s.borrowRepo.RejectPendingByBookTx(tx, book.ID, b.ID, "卖家已同意其他同学的借阅申请"); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}

		book.Status = constants.BookStatusLoaned
		if err := s.bookRepo.UpdateTx(tx, book); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}

		full, err := s.borrowRepo.FindByIDTx(tx, b.ID)
		if err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		d := s.toDTO(full)
		result = &d
		s.logger.Info(constants.LogBorrowApproved, "id", b.ID, "book_id", book.ID, "borrower", b.BorrowerID, "due", due.Format("2006-01-02"))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Reject lets the seller decline a pending request.
func (s *BorrowService) Reject(lenderID, borrowID uint, req dto.RejectBorrowRequest) (*dto.BorrowDTO, error) {
	b, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
	}
	if b.LenderID != lenderID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
	}
	if !b.IsPending() {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, constants.MsgBorrowStatusInvalid)
	}
	b.Status = constants.BorrowStatusRejected
	reason := req.Reason
	if reason == "" {
		reason = "卖家拒绝了借阅申请"
	}
	b.RejectReason = reason
	if err := s.borrowRepo.Update(b); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogBorrowRejected, "id", borrowID, "book_id", b.BookID, "operator", lenderID)
	d := s.toDTO(b)
	return &d, nil
}

// Return is the borrower action: I have given the book back, waiting for the
// seller to confirm. Allowed both before and after the due date.
func (s *BorrowService) Return(borrowerID, borrowID uint) (*dto.BorrowDTO, error) {
	var result *dto.BorrowDTO
	err := s.db.Transaction(func(tx *gorm.DB) error {
		b, err := s.borrowRepo.FindByIDForUpdate(tx, borrowID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
		}
		if b.BorrowerID != borrowerID {
			return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotBorrower)
		}
		if b.Status != constants.BorrowStatusApproved {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, constants.MsgBorrowStatusInvalid)
		}
		now := time.Now()
		b.Status = constants.BorrowStatusReturning
		b.ReturnedAt = &now
		if err := s.borrowRepo.UpdateTx(tx, b); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		full, err := s.borrowRepo.FindByIDTx(tx, b.ID)
		if err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		d := s.toDTO(full)
		result = &d
		s.logger.Info(constants.LogBorrowReturned, "id", borrowID, "book_id", b.BookID, "borrower", borrowerID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ConfirmReturn is the seller action after getting the book back: the borrow
// is closed and the book becomes borrowable (on_sale) again.
func (s *BorrowService) ConfirmReturn(lenderID, borrowID uint) (*dto.BorrowDTO, error) {
	var result *dto.BorrowDTO
	err := s.db.Transaction(func(tx *gorm.DB) error {
		b, err := s.borrowRepo.FindByIDForUpdate(tx, borrowID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
		}
		if b.LenderID != lenderID {
			return util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
		}
		if b.Status != constants.BorrowStatusReturning {
			return util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, constants.MsgBorrowStatusInvalid+": 借阅人尚未发起归还")
		}
		book, err := s.bookRepo.FindByIDForUpdate(tx, b.BookID)
		if err != nil {
			return util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound)
		}
		now := time.Now()
		b.Status = constants.BorrowStatusReturned
		b.ConfirmedAt = &now
		if err := s.borrowRepo.UpdateTx(tx, b); err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		// 恢复可借：借出中的书籍回到在售（可借标记保持不变）
		if book.Status == constants.BookStatusLoaned {
			book.Status = constants.BookStatusOnSale
			if err := s.bookRepo.UpdateTx(tx, book); err != nil {
				return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
			}
		}
		full, err := s.borrowRepo.FindByIDTx(tx, b.ID)
		if err != nil {
			return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		d := s.toDTO(full)
		result = &d
		s.logger.Info(constants.LogBorrowConfirmed, "id", borrowID, "book_id", b.BookID, "operator", lenderID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Remind lets the seller inspect the borrower and send an overdue reminder.
// Reminders go through the existing conversation/message channel and are only
// allowed while the book is actually lent out and past its due date.
func (s *BorrowService) Remind(lenderID, borrowID uint, req dto.RemindBorrowRequest) (*dto.BorrowDTO, error) {
	b, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
	}
	if b.LenderID != lenderID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, constants.MsgBorrowNotSeller)
	}
	now := time.Now()
	if !b.IsOverdue(now) {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, "仅逾期未还的借阅可以发送提醒")
	}
	if b.RemindedAt != nil && now.Sub(*b.RemindedAt) < reminderInterval {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeBorrowConflict, "提醒发送过于频繁，请稍后再试")
	}
	book, err := s.bookRepo.FindByID(b.BookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound)
	}
	content := req.Content
	if content == "" {
		content = defaultReminderContent(book.Title, *b.DueAt)
	}

	// 复用会话/消息：借阅沟通买家=借阅人，卖家=出借人，按书定位已有会话。
	conv, err := s.ensureConversation(b)
	if err != nil {
		return nil, err
	}
	msg := &model.Message{
		ConversationID: conv.ID,
		SenderID:       lenderID,
		Content:        content,
		IsRead:         false,
	}
	if err := s.msgRepo.Create(msg); err != nil {
		s.logger.Error("send reminder message failed", "borrow_id", borrowID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	_ = s.convRepo.Touch(conv.ID, truncate(content, 50))

	b.RemindedAt = &now
	if err := s.borrowRepo.Update(b); err != nil {
		s.logger.Warn("update reminded_at failed", "borrow_id", borrowID, "error", err)
	}
	s.logger.Info(constants.LogBorrowReminder, "id", borrowID, "book_id", b.BookID, "seller", lenderID, "borrower", b.BorrowerID)
	d := s.toDTO(b)
	return &d, nil
}

// List returns borrow requests visible to the user as lender or borrower.
func (s *BorrowService) List(q dto.BorrowQuery) (*dto.PageData, error) {
	if q.Role != "lender" && q.Role != "borrower" {
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "role 参数必须为 lender 或 borrower")
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = util.DefaultPageSize
	}
	if q.Size > util.MaxPageSize {
		q.Size = util.MaxPageSize
	}
	if q.Status != "" && !contains(constants.BorrowStatusOptions, q.Status) {
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeValidationError, "status 参数不合法")
	}
	items, total, err := s.borrowRepo.List(q)
	if err != nil {
		s.logger.Error(constants.LogBorrowList, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.BorrowDTO, 0, len(items))
	now := time.Now()
	for i := range items {
		list = append(list, s.toDTOWithNow(&items[i], now))
	}
	s.logger.Info(constants.LogBorrowList, "owner", q.Role, "user_id", q.LenderID+q.BorrowerID, "book_id", q.BookID)
	return &dto.PageData{List: list, Total: total, Page: q.Page, Size: q.Size}, nil
}

// Get returns a single borrow request if the caller is the lender or borrower.
func (s *BorrowService) Get(userID, borrowID uint) (*dto.BorrowDTO, error) {
	b, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
	}
	if b.LenderID != userID && b.BorrowerID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, "无权查看该借阅申请")
	}
	d := s.toDTO(b)
	return &d, nil
}

// Contact opens (or reuses) the conversation between lender and borrower for
// this borrow's book, so both sides can coordinate handover/return. Unlike
// CreateConversationFromBook it works for the seller too.
func (s *BorrowService) Contact(userID, borrowID uint, content string) (*dto.ConversationDTO, error) {
	b, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBorrowNotFound, constants.MsgNotFound+": borrow id="+fmt.Sprint(borrowID))
	}
	if b.LenderID != userID && b.BorrowerID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeBorrowForbidden, "无权查看该借阅申请")
	}
	conv, err := s.ensureConversation(b)
	if err != nil {
		return nil, err
	}
	if content != "" {
		msg := &model.Message{ConversationID: conv.ID, SenderID: userID, Content: content}
		if err := s.msgRepo.Create(msg); err != nil {
			return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
		}
		_ = s.convRepo.Touch(conv.ID, truncate(content, 50))
	}
	full, err := s.convRepo.FindByID(conv.ID)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	d := dto.FromConversation(full)
	return &d, nil
}

// ensureConversation returns the lender/borrower conversation for the book,
// creating it on first use.
func (s *BorrowService) ensureConversation(b *model.Borrow) (*model.Conversation, error) {
	conv, err := s.convRepo.FindExisting(b.BookID, 0, b.BorrowerID, b.LenderID)
	if err == nil {
		return conv, nil
	}
	conv = &model.Conversation{
		BookID:   b.BookID,
		BuyerID:  b.BorrowerID,
		SellerID: b.LenderID,
	}
	if err := s.convRepo.Create(conv); err != nil {
		s.logger.Error("create borrow conversation failed", "borrow_id", b.ID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	return conv, nil
}

func (s *BorrowService) toDTO(b *model.Borrow) dto.BorrowDTO {
	return s.toDTOWithNow(b, time.Now())
}

func (s *BorrowService) toDTOWithNow(b *model.Borrow, now time.Time) dto.BorrowDTO {
	d := dto.FromBorrow(b, now)
	fillBorrowRelations(&d, b)
	return d
}

func contains(options []string, v string) bool {
	for _, o := range options {
		if o == v {
			return true
		}
	}
	return false
}
