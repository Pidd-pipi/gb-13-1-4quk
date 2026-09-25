package service

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func seedUsers(t *testing.T, env *testEnv) (uint, uint) {
	t.Helper()
	repo := repository.NewUserRepository(env.db)
	u1 := &model.User{StudentNo: "S1001", Email: "s1@c.local", PasswordHash: "h", Department: "计算机学院", Role: constants.RoleStudent}
	u2 := &model.User{StudentNo: "S1002", Email: "s2@c.local", PasswordHash: "h", Department: "计算机学院", Role: constants.RoleStudent}
	if err := repo.Create(u1); err != nil {
		t.Fatalf("seed u1: %v", err)
	}
	if err := repo.Create(u2); err != nil {
		t.Fatalf("seed u2: %v", err)
	}
	return u1.ID, u2.ID
}

func newBookService(env *testEnv) *BookService {
	return NewBookService(
		env.db,
		repository.NewBookRepository(env.db),
		repository.NewFavoriteRepository(env.db),
		repository.NewBrowseHistoryRepository(env.db),
		repository.NewUserRepository(env.db),
		env.logger,
	)
}

func TestBookServiceLifecycle(t *testing.T) {
	env := newTestEnv(t)
	seller, buyer := seedUsers(t, env)
	svc := newBookService(env)

	req := dto.CreateBookRequest{
		Title: "高等数学", Author: "同济", ISBN: "978-1", Price: 20,
		Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience,
		TradeType: constants.TradeTypeInPerson, Campus: "东校区",
	}
	created, err := svc.CreateBook(seller, req)
	if err != nil {
		t.Fatalf("create book: %v", err)
	}
	if created.Status != constants.BookStatusOnSale {
		t.Errorf("status=%q, want on_sale", created.Status)
	}

	// list
	list, err := svc.ListBooks(dto.BookQuery{Keyword: "高等数学", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("list total=%d, want 1", list.Total)
	}

	// detail records view + history for buyer
	detail, err := svc.GetBookDetail(buyer, created.ID)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if detail.ViewCount <= 0 {
		t.Errorf("view_count=%d, want >0", detail.ViewCount)
	}

	// buyer reserves
	reserved, err := svc.ReserveBook(buyer, created.ID)
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if reserved.Status != constants.BookStatusReserved {
		t.Errorf("status=%q, want reserved", reserved.Status)
	}
	if reserved.ReservedBy != buyer {
		t.Errorf("reserved_by=%d, want %d", reserved.ReservedBy, buyer)
	}

	// buyer cannot reserve again (conflict)
	if _, err := svc.ReserveBook(buyer, created.ID); err == nil {
		t.Errorf("double reserve should fail")
	}

	// seller marks sold
	sold, err := svc.MarkSold(seller, created.ID)
	if err != nil {
		t.Fatalf("mark sold: %v", err)
	}
	if sold.Status != constants.BookStatusSold {
		t.Errorf("status=%q, want sold", sold.Status)
	}

	// sold is terminal
	if _, err := svc.MarkSold(seller, created.ID); err == nil {
		t.Errorf("sold -> sold should fail")
	}
	if _, err := svc.CancelReserve(buyer, created.ID); err == nil {
		t.Errorf("cancel on sold should fail")
	}
}

func TestBookServiceFavoritesAndRecommendations(t *testing.T) {
	env := newTestEnv(t)
	seller, buyer := seedUsers(t, env)
	svc := newBookService(env)

	created, err := svc.CreateBook(seller, dto.CreateBookRequest{
		Title: "数据结构", Price: 15, Condition: constants.ConditionSevenNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeMail,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := svc.AddFavorite(buyer, created.ID); err != nil {
		t.Fatalf("add favorite: %v", err)
	}
	if err := svc.AddFavorite(buyer, created.ID); err == nil {
		t.Errorf("duplicate favorite should fail")
	}
	favs, err := svc.ListFavorites(buyer, 1, 10)
	if err != nil {
		t.Fatalf("list favorites: %v", err)
	}
	if favs.Total != 1 {
		t.Errorf("favorites total=%d, want 1", favs.Total)
	}
	if err := svc.RemoveFavorite(buyer, created.ID); err != nil {
		t.Fatalf("remove favorite: %v", err)
	}
	if err := svc.RemoveFavorite(buyer, created.ID); err == nil {
		t.Errorf("double remove should fail")
	}

	recs, err := svc.GetRecommendations(buyer, 10)
	if err != nil {
		t.Fatalf("recommendations: %v", err)
	}
	if len(recs) != 1 {
		t.Errorf("recommendations len=%d, want 1 (same department)", len(recs))
	}
}

func TestBookServiceOwnershipGuard(t *testing.T) {
	env := newTestEnv(t)
	seller, other := seedUsers(t, env)
	svc := newBookService(env)
	created, err := svc.CreateBook(seller, dto.CreateBookRequest{
		Title: "线性代数", Price: 30, Condition: constants.ConditionBrandNew,
		SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.UpdateBook(other, created.ID, dto.UpdateBookRequest{Price: 1}); err == nil {
		t.Errorf("updating another user's book should fail")
	}
	if err := svc.DeleteBook(other, created.ID); err == nil {
		t.Errorf("deleting another user's book should fail")
	}
	if err := svc.DeleteBook(seller, created.ID); err != nil {
		t.Errorf("seller delete should succeed: %v", err)
	}
}
