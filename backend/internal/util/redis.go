package util

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/redis/go-redis/v9"
)

// CodeStore stores short-lived verification codes.
type CodeStore interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type redisCodeStore struct {
	client *redis.Client
}

func (s *redisCodeStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *redisCodeStore) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *redisCodeStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

type memoryCodeStore struct {
	mu    sync.Mutex
	items map[string]memoryItem
}

type memoryItem struct {
	value   string
	expires time.Time
}

func (s *memoryCodeStore) Set(_ context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = memoryItem{value: value, expires: time.Now().Add(ttl)}
	return nil
}

func (s *memoryCodeStore) Get(_ context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[key]
	if !ok {
		return "", redis.Nil
	}
	if time.Now().After(item.expires) {
		delete(s.items, key)
		return "", redis.Nil
	}
	return item.value, nil
}

func (s *memoryCodeStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
	return nil
}

// NewCodeStore builds a Redis-backed CodeStore, falling back to an in-memory
// implementation when Redis is unreachable so the platform still boots.
func NewCodeStore(cfg *config.Config, logger *slog.Logger) CodeStore {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn(constants.LogRedisConnectFailed, "error", err)
		return &memoryCodeStore{items: make(map[string]memoryItem)}
	}
	return &redisCodeStore{client: client}
}
