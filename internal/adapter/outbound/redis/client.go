package redis

import (
	"context"
	"crisplite/internal/domain"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func NewClient(ctx context.Context, redisCfg domain.RedisConfig) (*redis.Client, error) {
	var parseErr error
	once.Do(func() {
		opts, err := redis.ParseURL(redisCfg.URI)
		if err != nil {
			parseErr = err
			return
		}
		opts.PoolSize = 100
		redisClient = redis.NewClient(opts)
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
