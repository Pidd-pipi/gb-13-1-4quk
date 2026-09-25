package repository

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
)

func TestBookRepositoryListAndStatus(t *testing.T) {
	db := newTestDB(t)
	repo := NewBookRepository(db)

	// two sellers
	if err := db.Create(&model.User{StudentNo: "B001", Email: "b1@c.local", PasswordHash: "h", Role: "student"}).Error; err != nil {
		t.Fatalf("seed seller1: %v", err)
	}
	if err := db.Create(&model.User{StudentNo: "B002", Email: "b2@c.local", PasswordHash: "h", Role: "student"}).Error; err != nil {
		t.Fatalf("seed seller2: %v", err)
	}

	b1 := &model.Book{SellerID: 1, Title: "高等数学", Author: "同济", ISBN: "978-1", Price: 20, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Status: constants.BookStatusOnSale}
	if err := repo.Create(b1); err != nil {
		t.Fatalf("create b1: %v", err)
	}
	b2 := &model.Book{SellerID: 2, Title: "经济学原理", Author: "曼昆", ISBN: "978-2", Price: 30, Condition: constants.ConditionBrandNew, SubjectCategory: constants.SubjectEconManagement, TradeType: constants.TradeTypeMail, Status: constants.BookStatusOnSale}
	if err := repo.Create(b2); err != nil {
		t.Fatalf("create b2: %v", err)
	}

	q := dto.BookQuery{Keyword: "高等", Page: 1, PageSize: 10}
	items, total, err := repo.List(q)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Errorf("keyword filter total=%d len=%d, want 1/1", total, len(items))
	}
	if items[0].ID != b1.ID {
		t.Errorf("keyword filter returned wrong book")
	}

	q2 := dto.BookQuery{SubjectCategory: constants.SubjectScience, Sort: "price_asc", Page: 1, PageSize: 10}
	items2, total2, err := repo.List(q2)
	if err != nil {
		t.Fatalf("list2: %v", err)
	}
	if total2 != 1 {
		t.Errorf("category filter total=%d, want 1", total2)
	}
	_ = items2

	// for update + update
	tx := db.Begin()
	locked, err := repo.FindByIDForUpdate(tx, b1.ID)
	if err != nil {
		t.Fatalf("find for update: %v", err)
	}
	locked.Status = constants.BookStatusReserved
	if err := repo.UpdateTx(tx, locked); err != nil {
		t.Fatalf("update tx: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit: %v", err)
	}
	got, err := repo.FindByID(b1.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Status != constants.BookStatusReserved {
		t.Errorf("status = %q, want reserved", got.Status)
	}

	if err := repo.Delete(b2.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(b2.ID); err != ErrNotFound {
		t.Errorf("deleted book err = %v, want ErrNotFound", err)
	}
}

func TestBookRepositoryRecommendations(t *testing.T) {
	db := newTestDB(t)
	repo := NewBookRepository(db)
	_ = db.Create(&model.User{StudentNo: "C001", Email: "c1@c.local", PasswordHash: "h", Department: "计算机学院", Role: "student"}).Error
	_ = db.Create(&model.User{StudentNo: "C002", Email: "c2@c.local", PasswordHash: "h", Department: "计算机学院", Role: "student"}).Error
	_ = db.Create(&model.Book{SellerID: 1, Title: "A书", Price: 10, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Status: constants.BookStatusOnSale}).Error
	_ = db.Create(&model.Book{SellerID: 2, Title: "B书", Price: 12, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Status: constants.BookStatusOnSale}).Error
	items, err := repo.ListRecommendations("计算机学院", 1, 10)
	if err != nil {
		t.Fatalf("recommendations: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("recommendations len=%d, want 1", len(items))
	}
	if items[0].SellerID != 2 {
		t.Errorf("recommendation seller=%d, want 2", items[0].SellerID)
	}
}
