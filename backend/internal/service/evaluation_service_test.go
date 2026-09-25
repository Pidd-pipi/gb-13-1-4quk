package service

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func TestEvaluationService(t *testing.T) {
	env := newTestEnv(t)
	seller, buyer := seedUsers(t, env)
	bookRepo := repository.NewBookRepository(env.db)
	book := &model.Book{SellerID: seller, Title: "微观经济学", Price: 25, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectEconManagement, TradeType: constants.TradeTypeInPerson,
		Status: constants.BookStatusSold, ReservedBy: buyer}
	if err := bookRepo.Create(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	svc := NewEvaluationService(
		repository.NewEvaluationRepository(env.db),
		bookRepo,
		repository.NewUserRepository(env.db),
		env.logger,
	)
	eval, err := svc.CreateEvaluation(buyer, dto.CreateEvaluationRequest{
		ToUserID: seller, BookID: book.ID, Type: constants.EvaluationGood, Content: "书很好",
	})
	if err != nil {
		t.Fatalf("create evaluation: %v", err)
	}
	if eval.Type != constants.EvaluationGood {
		t.Errorf("type=%q, want good", eval.Type)
	}
	// duplicate evaluation fails
	if _, err := svc.CreateEvaluation(buyer, dto.CreateEvaluationRequest{
		ToUserID: seller, BookID: book.ID, Type: constants.EvaluationGood, Content: "再来一次",
	}); err == nil {
		t.Errorf("duplicate evaluation should fail")
	}
	list, err := svc.ListEvaluations(seller, 1, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("total=%d, want 1", list.Total)
	}
}

func TestEvaluationServiceOnlyAfterSold(t *testing.T) {
	env := newTestEnv(t)
	seller, buyer := seedUsers(t, env)
	bookRepo := repository.NewBookRepository(env.db)
	book := &model.Book{SellerID: seller, Title: "在售书", Price: 10, Condition: constants.ConditionNineNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Status: constants.BookStatusOnSale}
	_ = bookRepo.Create(book)
	svc := NewEvaluationService(
		repository.NewEvaluationRepository(env.db),
		bookRepo,
		repository.NewUserRepository(env.db),
		env.logger,
	)
	if _, err := svc.CreateEvaluation(buyer, dto.CreateEvaluationRequest{
		ToUserID: seller, BookID: book.ID, Type: constants.EvaluationGood,
	}); err == nil {
		t.Errorf("evaluation before sold should fail")
	}
}
