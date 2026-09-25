package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/util"
)

// AuthService handles registration, email verification and login.
type AuthService struct {
	userRepo  *repository.UserRepository
	codeStore util.CodeStore
	logger    *slog.Logger
	cfg       *config.Config
}

// NewAuthService creates an AuthService.
func NewAuthService(userRepo *repository.UserRepository, codeStore util.CodeStore, logger *slog.Logger, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, codeStore: codeStore, logger: logger, cfg: cfg}
}

// SendCode issues a 6-digit verification code for an email.
func (s *AuthService) SendCode(ctx context.Context, email string) (*dto.SendCodeResponse, error) {
	code, err := generateCode()
	if err != nil {
		s.logger.Error(constants.LogAuthCodeSent, "email", email, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	ttl := 10 * time.Minute
	if err := s.codeStore.Set(ctx, codeKey(email), code, ttl); err != nil {
		s.logger.Error("store email code failed", "email", email, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogAuthCodeSent, "email", email)
	return &dto.SendCodeResponse{Email: email, Expire: int(ttl.Seconds()), DevCode: code}, nil
}

// Register creates a student account after email verification.
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.TokenResponse, error) {
	stored, err := s.codeStore.Get(ctx, codeKey(req.Email))
	if err != nil {
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeEmailCodeInvalid, constants.MsgEmailCodeInvalid)
	}
	if stored != req.Code {
		s.logger.Warn(constants.LogAuthRegisterFailed, "student_no", req.StudentNo, "error", "invalid code")
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeEmailCodeInvalid, constants.MsgEmailCodeInvalid)
	}
	if _, err := s.codeStore.Get(ctx, usedCodeKey(req.Email)); err == nil {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeEmailTaken, "该邮箱已完成注册")
	}
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeEmailTaken, constants.MsgEmailTaken)
	}
	if _, err := s.userRepo.FindByStudentNo(req.StudentNo); err == nil {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeStudentNoTaken, constants.MsgStudentNoTaken)
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		s.logger.Error(constants.LogAuthRegisterFailed, "student_no", req.StudentNo, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	user := &model.User{
		StudentNo:     req.StudentNo,
		Email:         req.Email,
		PasswordHash:  hash,
		Role:          constants.RoleStudent,
		EmailVerified: true,
		Status:        constants.UserStatusActive,
	}
	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.NewAppError(http.StatusConflict, constants.CodeConflict, "学号或邮箱已存在")
		}
		s.logger.Error(constants.LogAuthRegisterFailed, "student_no", req.StudentNo, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	_ = s.codeStore.Set(ctx, usedCodeKey(req.Email), "1", 365*24*time.Hour)
	_ = s.codeStore.Delete(ctx, codeKey(req.Email))
	s.logger.Info(constants.LogAuthRegisterSuccess, "student_no", req.StudentNo, "email", req.Email)
	return s.issueToken(user)
}

// Login authenticates by email or student number.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.FindByAccount(req.Account)
	if err != nil {
		s.logger.Warn(constants.LogAuthLoginFailed, "email", req.Account, "error", "user not found")
		return nil, util.NewAppError(http.StatusUnauthorized, constants.CodeInvalidCredentials, constants.MsgInvalidCredentials)
	}
	if !util.CheckPassword(user.PasswordHash, req.Password) {
		s.logger.Warn(constants.LogAuthLoginFailed, "email", req.Account, "error", "password mismatch")
		return nil, util.NewAppError(http.StatusUnauthorized, constants.CodeInvalidCredentials, constants.MsgInvalidCredentials)
	}
	if user.Status != constants.UserStatusActive {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeUserDisabled, constants.MsgUserDisabled)
	}
	s.logger.Info(constants.LogAuthLoginSuccess, "email", req.Account, "role", user.Role)
	return s.issueToken(user)
}

func (s *AuthService) issueToken(user *model.User) (*dto.TokenResponse, error) {
	token, err := util.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpire)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.TokenResponse{Token: token, User: dto.FromUser(user)}, nil
}

func codeKey(email string) string { return "campusbooks:code:" + email }

func usedCodeKey(email string) string { return "campusbooks:used:" + email }

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
