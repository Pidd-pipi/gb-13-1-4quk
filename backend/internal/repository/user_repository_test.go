package repository

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

func TestUserRepositoryCRUD(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	hash, _ := util.HashPassword("test123")
	u := &model.User{StudentNo: "20231001", Email: "t1@campus.local", PasswordHash: hash, Role: "student"}
	if err := repo.Create(u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	got, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if got.StudentNo != "20231001" {
		t.Errorf("student_no = %q, want 20231001", got.StudentNo)
	}

	byEmail, err := repo.FindByEmail("t1@campus.local")
	if err != nil {
		t.Fatalf("find by email: %v", err)
	}
	if byEmail.ID != u.ID {
		t.Errorf("find by email returned wrong user")
	}

	byAccount, err := repo.FindByAccount("20231001")
	if err != nil {
		t.Fatalf("find by account: %v", err)
	}
	if byAccount.ID != u.ID {
		t.Errorf("find by account returned wrong user")
	}

	dupe := &model.User{StudentNo: "20231001", Email: "t2@campus.local", PasswordHash: hash, Role: "student"}
	if err := repo.Create(dupe); err != ErrDuplicate {
		t.Errorf("duplicate create err = %v, want ErrDuplicate", err)
	}
}

func TestUserRepositoryListAndStats(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	hash, _ := util.HashPassword("x")
	for i := 0; i < 3; i++ {
		u := &model.User{StudentNo: "S00" + string(rune('1'+i)), Email: "u" + string(rune('a'+i)) + "@c.local", PasswordHash: hash, Role: "student"}
		if err := repo.Create(u); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	items, total, err := repo.List("student", 0, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(items) != 2 {
		t.Errorf("len(items) = %d, want 2", len(items))
	}
	_, err = repo.FindByID(99999)
	if err != ErrNotFound {
		t.Errorf("missing user err = %v, want ErrNotFound", err)
	}
}
