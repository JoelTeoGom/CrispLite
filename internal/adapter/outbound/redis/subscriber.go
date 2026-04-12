package redis

import goredis "github.com/redis/go-redis/v9"

type Subscriber struct {
	client *goredis.Client
}
