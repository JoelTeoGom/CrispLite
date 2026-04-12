package redis

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStorage struct {
	client *redis.Client
	logger outbound.Logger
}

func NewCacheStorage(ctx context.Context, logger outbound.Logger, redisCfg domain.RedisConfig) (*CacheStorage, error) {
	client, err := NewClient(ctx, logger, redisCfg)
	if err != nil {
		return nil, err
	}
	return &CacheStorage{client: client, logger: logger}, nil
}

func (cs *CacheStorage) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := cs.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		cs.logger.ErrorWithVar(ctx, err, map[string]any{"key": key})
	}
	return err
}

func (cs *CacheStorage) Get(ctx context.Context, key string) (string, error) {
	result, err := cs.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		cs.logger.ErrorWithVar(ctx, err, map[string]any{"key": key})
	}
	return result, err
}

func (cs *CacheStorage) Delete(ctx context.Context, key string) error {
	err := cs.client.Del(ctx, key).Err()
	if err != nil {
		cs.logger.ErrorWithVar(ctx, err, map[string]any{"key": key})
	}
	return err
}

func (cs *CacheStorage) Close() error {
	return cs.client.Close()
}

func (cs *CacheStorage) Ping(ctx context.Context) error {
	err := cs.client.Ping(ctx).Err()
	if err != nil {
		cs.logger.Error(ctx, err)
	}
	return err
}

func (cs *CacheStorage) FlushAll(ctx context.Context) error {
	err := cs.client.FlushAll(ctx).Err()
	if err != nil {
		cs.logger.Error(ctx, err)
	}
	return err
}

func (cs *CacheStorage) Exists(ctx context.Context, key string) (bool, error) {
	count, err := cs.client.Exists(ctx, key).Result()
	if err != nil {
		cs.logger.ErrorWithVar(ctx, err, map[string]any{"key": key})
		return false, err
	}
	return count > 0, nil
}
