package service

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func TestConversationServiceBookFlow(t *testing.T) {
	env := newTestEnv(t)
	seller, buyer := seedUsers(t, env)
	bookRepo := repository.NewBookRepository(env.db)
	book := &model.Book{SellerID: seller, Title: "高等数学", Price: 20, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Status: constants.BookStatusOnSale}
	if err := bookRepo.Create(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	svc := NewConversationService(
		repository.NewConversationRepository(env.db),
		repository.NewMessageRepository(env.db),
		bookRepo,
		repository.NewWishRepository(env.db),
		env.logger,
	)
	conv, err := svc.CreateConversationFromBook(buyer, book.ID, "这本书还在吗？")
	if err != nil {
		t.Fatalf("create conversation: %v", err)
	}
	if conv.SellerID != seller || conv.BuyerID != buyer {
		t.Errorf("conv participants wrong: seller=%d buyer=%d", conv.SellerID, conv.BuyerID)
	}

	// same book -> reuse conversation, second message
	conv2, err := svc.CreateConversationFromBook(buyer, book.ID, "可以便宜一点吗？")
	if err != nil {
		t.Fatalf("create conversation again: %v", err)
	}
	if conv2.ID != conv.ID {
		t.Errorf("expected existing conversation reuse, got new id %d != %d", conv2.ID, conv.ID)
	}

	msgs, err := svc.ListMessages(buyer, conv.ID)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("messages len=%d, want 2", len(msgs))
	}

	// seller replies
	reply, err := svc.SendMessage(conv.ID, seller, dto.SendMessageRequest{Content: "可以，面交吧"})
	if err != nil {
		t.Fatalf("send message: %v", err)
	}
	if reply.SenderID != seller {
		t.Errorf("reply sender=%d, want %d", reply.SenderID, seller)
	}

	if err := svc.MarkRead(buyer, conv.ID); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	var unread int64
	if err := env.db.Model(&model.Message{}).Where("is_read = ? AND sender_id <> ?", false, buyer).Count(&unread).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if unread != 0 {
		t.Errorf("unread=%d, want 0", unread)
	}

	list, err := svc.ListConversations(buyer, 1, 10)
	if err != nil {
		t.Fatalf("list conversations: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("conversations total=%d, want 1", list.Total)
	}
}

func TestConversationServicePermission(t *testing.T) {
	env := newTestEnv(t)
	seller, buyer := seedUsers(t, env)
	bookRepo := repository.NewBookRepository(env.db)
	book := &model.Book{SellerID: seller, Title: "数据结构", Price: 15, Condition: constants.ConditionSevenNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeMail, Status: constants.BookStatusOnSale}
	_ = bookRepo.Create(book)

	svc := NewConversationService(
		repository.NewConversationRepository(env.db),
		repository.NewMessageRepository(env.db),
		bookRepo,
		repository.NewWishRepository(env.db),
		env.logger,
	)
	// seller cannot open conversation on own book
	if _, err := svc.CreateConversationFromBook(seller, book.ID, "hi"); err == nil {
		t.Errorf("seller contacting own book should fail")
	}
	conv, err := svc.CreateConversationFromBook(buyer, book.ID, "hi")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	// outsider cannot read messages
	outsider := &model.User{StudentNo: "O001", Email: "o@c.local", PasswordHash: "h", Role: constants.RoleStudent}
	_ = repository.NewUserRepository(env.db).Create(outsider)
	if _, err := svc.ListMessages(outsider.ID, conv.ID); err == nil {
		t.Errorf("outsider reading messages should fail")
	}
}

func TestConversationServiceSellerInitiated(t *testing.T) {
	env := newTestEnv(t)
	seller, borrower := seedUsers(t, env)
	bookRepo := repository.NewBookRepository(env.db)
	book := &model.Book{SellerID: seller, Title: "大学物理", Price: 15, Condition: constants.ConditionSevenNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Status: constants.BookStatusLentOut}
	if err := bookRepo.Create(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}
	svc := NewConversationService(
		repository.NewConversationRepository(env.db),
		repository.NewMessageRepository(env.db),
		bookRepo,
		repository.NewWishRepository(env.db),
		env.logger,
	)

	// non-seller cannot use the seller-initiated path
	if _, err := svc.CreateConversationToUser(borrower, book.ID, seller, "hi"); err == nil {
		t.Errorf("non-seller using seller-initiated path should fail")
	}
	// seller contacts the borrower about the lent-out book
	conv, err := svc.CreateConversationToUser(seller, book.ID, borrower, "同学，书快到期了")
	if err != nil {
		t.Fatalf("seller-initiated conversation: %v", err)
	}
	if conv.SellerID != seller || conv.BuyerID != borrower {
		t.Errorf("conv participants wrong: seller=%d buyer=%d", conv.SellerID, conv.BuyerID)
	}
	// SendBookMessage reuses the same conversation
	if err := svc.SendBookMessage(book.ID, seller, borrower, "【还书提醒】请尽快归还"); err != nil {
		t.Fatalf("send book message: %v", err)
	}
	msgs, err := svc.ListMessages(borrower, conv.ID)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("messages len=%d, want 2 (initial + reminder)", len(msgs))
	}
}
