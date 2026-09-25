package service

import (
	"context"
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func TestAuthServiceRegisterAndLogin(t *testing.T) {
	env := newTestEnv(t)
	userRepo := repository.NewUserRepository(env.db)
	svc := NewAuthService(userRepo, env.codeStore, env.logger, env.cfg)
	ctx := context.Background()

	codeResp, err := svc.SendCode(ctx, "new@campus.local")
	if err != nil {
		t.Fatalf("send code: %v", err)
	}
	if len(codeResp.DevCode) != 6 {
		t.Fatalf("dev_code len=%d, want 6", len(codeResp.DevCode))
	}

	token, err := svc.Register(ctx, dto.RegisterRequest{
		StudentNo: "20239999", Email: "new@campus.local", Password: "pass123", Code: codeResp.DevCode,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if token.Token == "" {
		t.Errorf("token empty")
	}
	if token.User.Role != constants.RoleStudent {
		t.Errorf("role=%q, want student", token.User.Role)
	}

	// duplicate email fails
	if _, err := svc.Register(ctx, dto.RegisterRequest{
		StudentNo: "20238888", Email: "new@campus.local", Password: "pass123", Code: codeResp.DevCode,
	}); err == nil {
		t.Errorf("duplicate email register should fail")
	}

	// login with email
	login, err := svc.Login(ctx, dto.LoginRequest{Account: "new@campus.local", Password: "pass123"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if login.User.Email != "new@campus.local" {
		t.Errorf("login email mismatch")
	}
	// login with student no
	if _, err := svc.Login(ctx, dto.LoginRequest{Account: "20239999", Password: "pass123"}); err != nil {
		t.Fatalf("login by student no: %v", err)
	}
	// wrong password
	if _, err := svc.Login(ctx, dto.LoginRequest{Account: "new@campus.local", Password: "wrong"}); err == nil {
		t.Errorf("wrong password should fail")
	}
}

func TestAuthServiceInvalidCode(t *testing.T) {
	env := newTestEnv(t)
	svc := NewAuthService(repository.NewUserRepository(env.db), env.codeStore, env.logger, env.cfg)
	_, err := svc.Register(context.Background(), dto.RegisterRequest{
		StudentNo: "20237777", Email: "x@campus.local", Password: "pass123", Code: "000000",
	})
	if err == nil {
		t.Errorf("register with invalid code should fail")
	}
}
