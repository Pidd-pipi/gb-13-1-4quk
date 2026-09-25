package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("JWT_SECRET")
	cfg := Load()
	if cfg.ServerPort != "8080" {
		t.Errorf("ServerPort default = %q, want 8080", cfg.ServerPort)
	}
	if cfg.DBHost != "db" {
		t.Errorf("DBHost default = %q, want db", cfg.DBHost)
	}
	if cfg.JWTSecret != "change_me_to_a_long_random_string" {
		t.Errorf("JWTSecret default mismatch")
	}
	if cfg.RateLimitReq <= 0 {
		t.Errorf("RateLimitReq must be positive")
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "44007")
	t.Setenv("JWT_EXPIRE_HOURS", "24")
	cfg := Load()
	if cfg.ServerPort != "9090" {
		t.Errorf("ServerPort = %q, want 9090", cfg.ServerPort)
	}
	if cfg.DBHost != "127.0.0.1" {
		t.Errorf("DBHost = %q, want 127.0.0.1", cfg.DBHost)
	}
	if cfg.JWTExpire != 24*time.Hour {
		t.Errorf("JWTExpire = %v, want 24h", cfg.JWTExpire)
	}
	if got := cfg.DSN(); got == "" {
		t.Errorf("DSN must not be empty")
	}
}
