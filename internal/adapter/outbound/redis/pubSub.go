package redis

import goredis "github.com/redis/go-redis/v9"

type Publisher struct {
	client *goredis.Client
}
