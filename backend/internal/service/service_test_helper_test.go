package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type testEnv struct {
	db        *gorm.DB
	codeStore util.CodeStore
	logger    *slog.Logger
	cfg       *config.Config
}

// memStore is an in-memory CodeStore for tests.
type memStore struct {
	mu sync.Mutex
	m  map[string]memItem
}

type memItem struct {
	v   string
	exp time.Time
}

func (s *memStore) Set(_ context.Context, k, v string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k] = memItem{v: v, exp: time.Now().Add(ttl)}
	return nil
}

func (s *memStore) Get(_ context.Context, k string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.m[k]
	if !ok || time.Now().After(it.exp) {
		return "", errNotFound
	}
	return it.v, nil
}

func (s *memStore) Delete(_ context.Context, k string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, k)
	return nil
}

type notFoundError struct{}

func (*notFoundError) Error() string { return "not found" }

var errNotFound error = &notFoundError{}

// newTestEnv wires the service dependencies against a uniquely-named
// in-memory sqlite DB (unique name avoids parallel-package shared-cache issues).
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	dsn := fmt.Sprintf("file:srv_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(
		&model.User{}, &model.Book{}, &model.Wish{}, &model.Conversation{},
		&model.Message{}, &model.Evaluation{}, &model.Favorite{},
		&model.BrowseHistory{}, &model.AuditLog{}, &model.BorrowRequest{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	cfg := &config.Config{
		ServerPort: "8080", JWTSecret: "test-secret", JWTExpire: time.Hour,
		DBHost: "localhost", DBPort: "3306", DBUser: "u", DBPassword: "p", DBName: "n",
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return &testEnv{db: db, codeStore: &memStore{m: make(map[string]memItem)}, logger: logger, cfg: cfg}
}
