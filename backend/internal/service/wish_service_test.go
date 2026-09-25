package service

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func TestWishServiceLifecycle(t *testing.T) {
	env := newTestEnv(t)
	userRepo := repository.NewUserRepository(env.db)
	u := &model.User{StudentNo: "W001", Email: "w@c.local", PasswordHash: "h", Role: constants.RoleStudent}
	if err := userRepo.Create(u); err != nil {
		t.Fatalf("seed: %v", err)
	}
	svc := NewWishService(repository.NewWishRepository(env.db), env.logger)

	created, err := svc.CreateWish(u.ID, dto.CreateWishRequest{
		BookTitle: "操作系统概念", Author: "Silberschatz", ExpectedPrice: 30,
		ConditionRequirement: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience,
	})
	if err != nil {
		t.Fatalf("create wish: %v", err)
	}
	if created.Status != constants.WishStatusOpen {
		t.Errorf("status=%q, want open", created.Status)
	}

	list, err := svc.ListWishes(dto.WishQuery{Keyword: "操作系统", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("list total=%d, want 1", list.Total)
	}

	closed, err := svc.CloseWish(u.ID, created.ID)
	if err != nil {
		t.Fatalf("close: %v", err)
	}
	if closed.Status != constants.WishStatusClosed {
		t.Errorf("status=%q, want closed", closed.Status)
	}
	if _, err := svc.CloseWish(u.ID, created.ID); err == nil {
		t.Errorf("double close should fail")
	}
	if err := svc.DeleteWish(u.ID, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.GetWish(created.ID); err == nil {
		t.Errorf("deleted wish should be gone")
	}
}

func TestWishServiceOwnershipGuard(t *testing.T) {
	env := newTestEnv(t)
	userRepo := repository.NewUserRepository(env.db)
	u1 := &model.User{StudentNo: "W101", Email: "w1@c.local", PasswordHash: "h", Role: constants.RoleStudent}
	u2 := &model.User{StudentNo: "W102", Email: "w2@c.local", PasswordHash: "h", Role: constants.RoleStudent}
	_ = userRepo.Create(u1)
	_ = userRepo.Create(u2)
	svc := NewWishService(repository.NewWishRepository(env.db), env.logger)
	created, err := svc.CreateWish(u1.ID, dto.CreateWishRequest{
		BookTitle: "管理学", SubjectCategory: constants.SubjectEconManagement,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.CloseWish(u2.ID, created.ID); err == nil {
		t.Errorf("closing other's wish should fail")
	}
	if err := svc.DeleteWish(u2.ID, created.ID); err == nil {
		t.Errorf("deleting other's wish should fail")
	}
}
