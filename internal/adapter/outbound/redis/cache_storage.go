package redis

import (
	"context"
	locallogger "crisplite/internal/adapter/outbound/local_logger"
	"crisplite/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStorage struct {
	client *redis.Client
}

func NewCacheStorage(ctx context.Context, logger locallogger.LocalLogger, redisCfg domain.RedisConfig) (*CacheStorage, error) {
	client, err := NewClient(ctx, redisCfg)
	if err != nil {
		return nil, err
	}
	return &CacheStorage{client: client}, nil
}

func (cs *CacheStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return cs.client.Set(ctx, key, value, expiration).Err()
}

func (cs *CacheStorage) Get(ctx context.Context, key string) (string, error) {
	return cs.client.Get(ctx, key).Result()
}

func (cs *CacheStorage) Delete(ctx context.Context, key string) error {
	return cs.client.Del(ctx, key).Err()
}

func (cs *CacheStorage) Close() error {
	return cs.client.Close()
}

func (cs *CacheStorage) Ping(ctx context.Context) error {
	return cs.client.Ping(ctx).Err()
}

func (cs *CacheStorage) FlushAll(ctx context.Context) error {
	return cs.client.FlushAll(ctx).Err()
}

func (cs *CacheStorage) Exists(ctx context.Context, key string) (bool, error) {
	count, err := cs.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
