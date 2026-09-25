package service

import (
	"testing"
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func newBorrowService(env *testEnv) *BorrowService {
	bookRepo := repository.NewBookRepository(env.db)
	convService := NewConversationService(
		repository.NewConversationRepository(env.db),
		repository.NewMessageRepository(env.db),
		bookRepo,
		repository.NewWishRepository(env.db),
		env.logger,
	)
	return NewBorrowService(env.db, repository.NewBorrowRepository(env.db), bookRepo, convService, env.logger)
}

func seedThirdUser(t *testing.T, env *testEnv) uint {
	t.Helper()
	repo := repository.NewUserRepository(env.db)
	u := &model.User{StudentNo: "S1003", Email: "s3@c.local", PasswordHash: "h", Department: "计算机学院", Role: constants.RoleStudent}
	if err := repo.Create(u); err != nil {
		t.Fatalf("seed u3: %v", err)
	}
	return u.ID
}

func createLendableBook(t *testing.T, env *testEnv, sellerID uint, lendDays int) *dto.BookDTO {
	t.Helper()
	svc := newBookService(env)
	book, err := svc.CreateBook(sellerID, dto.CreateBookRequest{
		Title: "短借教材", Price: 10, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson,
		Lendable: true, LendDays: lendDays,
	})
	if err != nil {
		t.Fatalf("create lendable book: %v", err)
	}
	return book
}

func TestBorrowFlowLifecycle(t *testing.T) {
	env := newTestEnv(t)
	seller, borrower := seedUsers(t, env)
	svc := newBorrowService(env)
	bookSvc := newBookService(env)
	book := createLendableBook(t, env, seller, constants.LendDays14)

	// borrower applies
	req, err := svc.Apply(borrower, book.ID)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if req.Status != constants.BorrowStatusPending {
		t.Errorf("status=%q, want pending", req.Status)
	}

	// duplicate application blocked
	if _, err := svc.Apply(borrower, book.ID); err == nil {
		t.Errorf("duplicate apply should fail")
	}

	// borrower cannot approve (not the seller)
	if _, err := svc.Approve(borrower, req.ID); err == nil {
		t.Errorf("borrower approve should fail")
	}

	// seller approves -> book lent out with due date
	approved, err := svc.Approve(seller, req.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if approved.Status != constants.BorrowStatusApproved {
		t.Errorf("status=%q, want approved", approved.Status)
	}
	if approved.DueAt == "" {
		t.Errorf("due_at should be set after approval")
	}
	detail, err := bookSvc.GetBookDetail(seller, book.ID)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if detail.Status != constants.BookStatusLentOut {
		t.Errorf("book status=%q, want lent_out", detail.Status)
	}
	if detail.ActiveBorrow == nil || detail.ActiveBorrow.Borrower == nil {
		t.Errorf("seller should see active borrow with borrower info")
	}

	// lent-out book cannot be reserved / sold / deleted
	if _, err := bookSvc.ReserveBook(seller, book.ID); err == nil {
		t.Errorf("reserve on lent-out book should fail")
	}
	if _, err := bookSvc.MarkSold(seller, book.ID); err == nil {
		t.Errorf("mark sold on lent-out book should fail")
	}
	if err := bookSvc.DeleteBook(seller, book.ID); err == nil {
		t.Errorf("delete on lent-out book should fail")
	}

	// borrower returns
	returned, err := svc.Return(borrower, req.ID)
	if err != nil {
		t.Fatalf("return: %v", err)
	}
	if returned.Status != constants.BorrowStatusReturned {
		t.Errorf("status=%q, want returned", returned.Status)
	}

	// seller confirms -> book back on sale
	completed, err := svc.ConfirmReturn(seller, req.ID)
	if err != nil {
		t.Fatalf("confirm return: %v", err)
	}
	if completed.Status != constants.BorrowStatusCompleted {
		t.Errorf("status=%q, want completed", completed.Status)
	}
	detail, err = bookSvc.GetBookDetail(0, book.ID)
	if err != nil {
		t.Fatalf("detail after confirm: %v", err)
	}
	if detail.Status != constants.BookStatusOnSale {
		t.Errorf("book status=%q, want on_sale after return confirm", detail.Status)
	}
	if detail.ActiveBorrow != nil {
		t.Errorf("active borrow should be cleared after completion")
	}
}

func TestBorrowApproveRejectsOtherApplications(t *testing.T) {
	env := newTestEnv(t)
	seller, borrowerA := seedUsers(t, env)
	borrowerB := seedThirdUser(t, env)
	svc := newBorrowService(env)
	book := createLendableBook(t, env, seller, constants.LendDays7)

	reqA, err := svc.Apply(borrowerA, book.ID)
	if err != nil {
		t.Fatalf("apply A: %v", err)
	}
	reqB, err := svc.Apply(borrowerB, book.ID)
	if err != nil {
		t.Fatalf("apply B: %v", err)
	}

	if _, err := svc.Approve(seller, reqA.ID); err != nil {
		t.Fatalf("approve A: %v", err)
	}

	// the other pending application is auto-rejected
	reloaded, err := svc.ListMine(borrowerB, "borrower", 1, 10)
	if err != nil {
		t.Fatalf("list mine: %v", err)
	}
	items := reloaded.List.([]dto.BorrowDTO)
	if len(items) != 1 || items[0].ID != reqB.ID || items[0].Status != constants.BorrowStatusRejected {
		t.Errorf("borrower B request should be rejected, got %+v", items)
	}

	// no new applications while the book is lent out
	if _, err := svc.Apply(borrowerB, book.ID); err == nil {
		t.Errorf("apply on lent-out book should fail")
	}
}

func TestBorrowApplyGuards(t *testing.T) {
	env := newTestEnv(t)
	seller, borrower := seedUsers(t, env)
	svc := newBorrowService(env)
	bookSvc := newBookService(env)

	// cannot borrow own book
	own := createLendableBook(t, env, seller, constants.LendDays7)
	if _, err := svc.Apply(seller, own.ID); err == nil {
		t.Errorf("apply to own book should fail")
	}

	// cannot apply to a non-lendable book
	plain, err := bookSvc.CreateBook(seller, dto.CreateBookRequest{
		Title: "普通出售书", Price: 10, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson,
	})
	if err != nil {
		t.Fatalf("create plain book: %v", err)
	}
	if _, err := svc.Apply(borrower, plain.ID); err == nil {
		t.Errorf("apply to non-lendable book should fail")
	}

	// lendable flag requires 7/14 days
	if _, err := bookSvc.CreateBook(seller, dto.CreateBookRequest{
		Title: "非法借期", Price: 10, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson,
		Lendable: true, LendDays: 30,
	}); err == nil {
		t.Errorf("lendable book with invalid lend days should fail")
	}

	// reject flow: seller rejects, borrower may apply again
	req, err := svc.Apply(borrower, own.ID)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	rejected, err := svc.Reject(seller, req.ID)
	if err != nil {
		t.Fatalf("reject: %v", err)
	}
	if rejected.Status != constants.BorrowStatusRejected {
		t.Errorf("status=%q, want rejected", rejected.Status)
	}
	if _, err := svc.Apply(borrower, own.ID); err != nil {
		t.Errorf("re-apply after rejection should succeed: %v", err)
	}
}

func TestBorrowRemindAndOverdue(t *testing.T) {
	env := newTestEnv(t)
	seller, borrower := seedUsers(t, env)
	svc := newBorrowService(env)
	book := createLendableBook(t, env, seller, constants.LendDays7)

	req, err := svc.Apply(borrower, book.ID)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	// cannot remind before approval
	if _, err := svc.Remind(seller, req.ID); err == nil {
		t.Errorf("remind on pending request should fail")
	}

	if _, err := svc.Approve(seller, req.ID); err != nil {
		t.Fatalf("approve: %v", err)
	}

	// borrower cannot send reminder
	if _, err := svc.Remind(borrower, req.ID); err == nil {
		t.Errorf("borrower remind should fail")
	}

	// seller reminds -> message delivered via conversation, reminded_at recorded
	reminded, err := svc.Remind(seller, req.ID)
	if err != nil {
		t.Fatalf("remind: %v", err)
	}
	if reminded.RemindedAt == "" {
		t.Errorf("reminded_at should be set")
	}
	var convCount int64
	if err := env.db.Model(&model.Conversation{}).
		Where("book_id = ? AND buyer_id = ? AND seller_id = ?", book.ID, borrower, seller).
		Count(&convCount).Error; err != nil {
		t.Fatalf("count conversations: %v", err)
	}
	if convCount != 1 {
		t.Errorf("conversation count=%d, want 1", convCount)
	}
	var msgCount int64
	if err := env.db.Model(&model.Message{}).Count(&msgCount).Error; err != nil {
		t.Fatalf("count messages: %v", err)
	}
	if msgCount != 1 {
		t.Errorf("message count=%d, want 1", msgCount)
	}

	// force the due date into the past -> overdue shows up in listings
	past := time.Now().Add(-time.Hour)
	if err := env.db.Model(&model.BorrowRequest{}).Where("id = ?", req.ID).Update("due_at", past).Error; err != nil {
		t.Fatalf("backdate due_at: %v", err)
	}
	mine, err := svc.ListMine(borrower, "borrower", 1, 10)
	if err != nil {
		t.Fatalf("list mine: %v", err)
	}
	items := mine.List.([]dto.BorrowDTO)
	if len(items) != 1 || !items[0].Overdue {
		t.Errorf("request should be overdue after due date passed, got %+v", items)
	}
}
