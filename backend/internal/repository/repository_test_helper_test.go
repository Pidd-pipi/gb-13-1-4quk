package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// newTestDB builds a uniquely-named in-memory sqlite DB with all models migrated.
// A unique name avoids cross-package shared-cache interference when go test runs
// packages in parallel.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(
		&model.User{},
		&model.Book{},
		&model.Wish{},
		&model.Conversation{},
		&model.Message{},
		&model.Evaluation{},
		&model.Favorite{},
		&model.BrowseHistory{},
		&model.AuditLog{},
		&model.Borrow{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}
