package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Sentinel errors shared by all repositories.
var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate entry") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "Error 1062")
}

// translate maps GORM errors to sentinel errors.
func translate(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}
