package redis

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func NewClient(ctx context.Context, logger outbound.Logger, redisCfg domain.RedisConfig) (*redis.Client, error) {
	var parseErr error
	once.Do(func() {
		opts, err := redis.ParseURL(redisCfg.URI)
		if err != nil {
			logger.Error(ctx, fmt.Errorf("failed to parse redis URL: %w", err))
			parseErr = err
			return
		}
		opts.PoolSize = 100
		redisClient = redis.NewClient(opts)
		logger.Info(ctx, "redis client initialized")
	})

	if parseErr != nil {
		return nil, parseErr
	}
	return redisClient, nil
}

func GetClient() *redis.Client {
	return redisClient
}

func CloseClient() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
