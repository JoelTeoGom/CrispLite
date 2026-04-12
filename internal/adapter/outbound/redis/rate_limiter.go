package redis

import goredis "github.com/redis/go-redis/v9"

type RateLimiter struct {
	client *goredis.Client
}
