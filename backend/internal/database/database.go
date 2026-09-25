package database

import (
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/model"
)

// Connect opens the MySQL connection and retries until ready.
func Connect(cfg *config.Config, logger *slog.Logger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
			Logger:                                   gormlogger.Default.LogMode(gormlogger.Warn),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			var sqlDB interface{ Ping() error }
			_ = sqlDB
			if raw, e := db.DB(); e == nil {
				raw.SetMaxOpenConns(50)
				raw.SetMaxIdleConns(10)
				raw.SetConnMaxLifetime(5 * time.Minute)
				if e := raw.Ping(); e == nil {
					return db, nil
				}
			}
		}
		logger.Warn("database not ready, retrying", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

// Migrate creates/updates all tables via GORM AutoMigrate.
func Migrate(db *gorm.DB, logger *slog.Logger) error {
	logger.Info("running database migrations")
	return db.AutoMigrate(
		&model.User{},
		&model.Book{},
		&model.Wish{},
		&model.Conversation{},
		&model.Message{},
		&model.Evaluation{},
		&model.Favorite{},
		&model.BrowseHistory{},
		&model.AuditLog{},
	)
}
