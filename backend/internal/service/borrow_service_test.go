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
	return NewBorrowService(
		env.db,
		repository.NewBorrowRepository(env.db),
		repository.NewBookRepository(env.db),
		repository.NewConversationRepository(env.db),
		repository.NewMessageRepository(env.db),
		env.logger,
	)
}

func borrowableBook(t *testing.T, env *testEnv, svc *BookService, seller uint, duration int) *dto.BookDTO {
	t.Helper()
	book, err := svc.CreateBook(seller, dto.CreateBookRequest{
		Title: "借来过渡的教材", Price: 20, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson,
		Borrowable: true, BorrowDuration: duration,
	})
	if err != nil {
		t.Fatalf("create borrowable book: %v", err)
	}
	return book
}

func TestBorrowFlowApplyApproveReturnConfirm(t *testing.T) {
	env := newTestEnv(t)
	lender, borrower := seedUsers(t, env)
	bookSvc := newBookService(env)
	svc := newBorrowService(env)
	book := borrowableBook(t, env, bookSvc, lender, 7)

	// 1. 借阅人提交申请
	req := dto.CreateBorrowRequest{}
	b, err := svc.Apply(borrower, book.ID, req)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if b.Status != constants.BorrowStatusPending {
		t.Errorf("status=%q, want pending", b.Status)
	}

	// 2. 重复申请被拒
	if _, err := svc.Apply(borrower, book.ID, req); err == nil {
		t.Errorf("duplicate pending apply should fail")
	}

	// 3. 卖家不能借自己的书
	if _, err := svc.Apply(lender, book.ID, req); err == nil {
		t.Errorf("self borrow should fail")
	}

	// 4. 非卖家不能同意
	if _, err := svc.Approve(borrower, b.ID); err == nil {
		t.Errorf("non-owner approve should fail")
	}

	// 5. 卖家同意 -> 书籍借出、记录到期日
	approved, err := svc.Approve(lender, b.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if approved.Status != constants.BorrowStatusApproved || approved.DueAt == "" {
		t.Errorf("approved=%q due=%q, want approved+due date", approved.Status, approved.DueAt)
	}
	detail, err := bookSvc.GetBookDetail(borrower, book.ID)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if detail.Status != constants.BookStatusLoaned || detail.ActiveBorrow == nil {
		t.Errorf("book status=%q active=%v, want loaned+active borrow", detail.Status, detail.ActiveBorrow != nil)
	}
	if detail.ActiveBorrow.IsOverdue {
		t.Errorf("fresh loan should not be overdue")
	}

	// 6. 借出后其他申请被拒
	other := createBorrower(t, env, "S2001", "o1@c.local")
	if _, err := svc.Apply(other, book.ID, req); err == nil {
		t.Errorf("apply while loaned should fail")
	}

	// 7. 借阅人发起归还 -> returning
	returning, err := svc.Return(borrower, b.ID)
	if err != nil {
		t.Fatalf("return: %v", err)
	}
	if returning.Status != constants.BorrowStatusReturning {
		t.Errorf("status=%q, want returning", returning.Status)
	}
	// 非卖家不能确认收回
	if _, err := svc.ConfirmReturn(borrower, b.ID); err == nil {
		t.Errorf("borrower confirm-return should fail")
	}

	// 8. 卖家确认收回 -> returned，书籍恢复可借（on_sale）
	confirmed, err := svc.ConfirmReturn(lender, b.ID)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if confirmed.Status != constants.BorrowStatusReturned || confirmed.ConfirmedAt == "" {
		t.Errorf("status=%q confirmed_at=%q, want returned+confirmed", confirmed.Status, confirmed.ConfirmedAt)
	}
	detail2, err := bookSvc.GetBookDetail(borrower, book.ID)
	if err != nil {
		t.Fatalf("detail2: %v", err)
	}
	if detail2.Status != constants.BookStatusOnSale || !detail2.Borrowable {
		t.Errorf("after return status=%q borrowable=%v, want on_sale/borrowable", detail2.Status, detail2.Borrowable)
	}
	if detail2.ActiveBorrow != nil {
		t.Errorf("completed borrow should not be active")
	}
}

func TestBorrowApproveRejectsOtherPending(t *testing.T) {
	env := newTestEnv(t)
	lender, borrower := seedUsers(t, env)
	bookSvc := newBookService(env)
	svc := newBorrowService(env)
	book := borrowableBook(t, env, bookSvc, lender, 14)

	first, err := svc.Apply(borrower, book.ID, dto.CreateBorrowRequest{})
	if err != nil {
		t.Fatalf("apply first: %v", err)
	}
	other := createBorrower(t, env, "S2002", "o2@c.local")
	second, err := svc.Apply(other, book.ID, dto.CreateBorrowRequest{})
	if err != nil {
		t.Fatalf("apply second: %v", err)
	}
	if _, err := svc.Approve(lender, first.ID); err != nil {
		t.Fatalf("approve first: %v", err)
	}
	got, err := svc.Get(other, second.ID)
	if err != nil {
		t.Fatalf("get second: %v", err)
	}
	if got.Status != constants.BorrowStatusRejected {
		t.Errorf("other pending request status=%q, want rejected", got.Status)
	}
}

func TestBorrowOverdueAndReminder(t *testing.T) {
	env := newTestEnv(t)
	lender, borrower := seedUsers(t, env)
	bookSvc := newBookService(env)
	svc := newBorrowService(env)
	book := borrowableBook(t, env, bookSvc, lender, 7)
	b, err := svc.Apply(borrower, book.ID, dto.CreateBorrowRequest{})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if _, err := svc.Approve(lender, b.ID); err != nil {
		t.Fatalf("approve: %v", err)
	}

	// 未到期不能提醒
	if _, err := svc.Remind(lender, b.ID, dto.RemindBorrowRequest{}); err == nil {
		t.Errorf("reminder before due should fail")
	}

	// 把到期日挪到过去，模拟逾期
	repo := repository.NewBorrowRepository(env.db)
	stored, err := repo.FindByID(b.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	past := time.Now().Add(-24 * time.Hour)
	stored.DueAt = &past
	if err := repo.Update(stored); err != nil {
		t.Fatalf("backdate: %v", err)
	}

	got, err := svc.Get(borrower, b.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !got.IsOverdue {
		t.Errorf("loan past due should be overdue")
	}

	// 逾期后卖家发送提醒（复用站内消息）
	if _, err := svc.Remind(lender, b.ID, dto.RemindBorrowRequest{}); err != nil {
		t.Fatalf("remind: %v", err)
	}
	convRepo := repository.NewConversationRepository(env.db)
	msgRepo := repository.NewMessageRepository(env.db)
	convs, total, err := convRepo.ListByUser(lender, 0, 10)
	if err != nil {
		t.Fatalf("list convs: %v", err)
	}
	if total != 1 {
		t.Fatalf("conversations=%d, want 1 reminder conversation", total)
	}
	msgs, err := msgRepo.ListByConversation(convs[0].ID)
	if err != nil || len(msgs) != 1 {
		t.Fatalf("reminder messages=%d err=%v, want 1", len(msgs), err)
	}
	if msgs[0].SenderID != lender || msgs[0].Content == "" {
		t.Errorf("reminder message sender=%d content=%q", msgs[0].SenderID, msgs[0].Content)
	}

	// 非卖家不能提醒
	if _, err := svc.Remind(borrower, b.ID, dto.RemindBorrowRequest{}); err == nil {
		t.Errorf("borrower remind should fail")
	}
}

func TestBorrowNotAllowedAndReject(t *testing.T) {
	env := newTestEnv(t)
	lender, borrower := seedUsers(t, env)
	bookSvc := newBookService(env)
	svc := newBorrowService(env)

	// 未开启短借的书不能申请
	plain, err := bookSvc.CreateBook(lender, dto.CreateBookRequest{
		Title: "只卖不借", Price: 30, Condition: constants.ConditionBrandNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson,
	})
	if err != nil {
		t.Fatalf("create plain: %v", err)
	}
	if _, err := svc.Apply(borrower, plain.ID, dto.CreateBorrowRequest{}); err == nil {
		t.Errorf("apply on non-borrowable book should fail")
	}

	book := borrowableBook(t, env, bookSvc, lender, 7)
	b, err := svc.Apply(borrower, book.ID, dto.CreateBorrowRequest{})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	rejected, err := svc.Reject(lender, b.ID, dto.RejectBorrowRequest{Reason: "这本书我还要用"})
	if err != nil {
		t.Fatalf("reject: %v", err)
	}
	if rejected.Status != constants.BorrowStatusRejected || rejected.RejectReason == "" {
		t.Errorf("status=%q reason=%q", rejected.Status, rejected.RejectReason)
	}
	// 已拒绝的申请不能再同意
	if _, err := svc.Approve(lender, b.ID); err == nil {
		t.Errorf("approve rejected request should fail")
	}
}

func TestBorrowContactOpensConversationForBothSides(t *testing.T) {
	env := newTestEnv(t)
	lender, borrower := seedUsers(t, env)
	bookSvc := newBookService(env)
	svc := newBorrowService(env)
	book := borrowableBook(t, env, bookSvc, lender, 7)
	b, err := svc.Apply(borrower, book.ID, dto.CreateBorrowRequest{})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if _, err := svc.Approve(lender, b.ID); err != nil {
		t.Fatalf("approve: %v", err)
	}
	// 卖家也能通过借阅记录发起会话（普通书籍会话会拒绝卖家）
	if _, err := svc.Contact(lender, b.ID, "方便什么时候还书？"); err != nil {
		t.Fatalf("lender contact: %v", err)
	}
	// 复用同一会话
	if _, err := svc.Contact(borrower, b.ID, "明天下午可以吗？"); err != nil {
		t.Fatalf("borrower contact: %v", err)
	}
	convRepo := repository.NewConversationRepository(env.db)
	_, total, err := convRepo.ListByUser(lender, 0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 {
		t.Errorf("conversations=%d, want 1 shared", total)
	}
	// 无关用户不能联系
	other := createBorrower(t, env, "S3001", "x@c.local")
	if _, err := svc.Contact(other, b.ID, ""); err == nil {
		t.Errorf("outsider contact should fail")
	}
}

// createBorrower seeds an extra student beyond the default seller/buyer pair.
func createBorrower(t *testing.T, env *testEnv, studentNo, email string) uint {
	t.Helper()
	repo := repository.NewUserRepository(env.db)
	u := &model.User{StudentNo: studentNo, Email: email, PasswordHash: "h", Department: "计算机学院", Role: constants.RoleStudent}
	if err := repo.Create(u); err != nil {
		t.Fatalf("seed %s: %v", email, err)
	}
	return u.ID
}
