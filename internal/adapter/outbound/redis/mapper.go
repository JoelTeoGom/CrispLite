package redis

import (
	"crisplite/internal/domain"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func MapRedisMessageToDomainMessage(msg *redis.Message) (*domain.Message, error) {
	var dm domain.Message
	if err := json.Unmarshal([]byte(msg.Payload), &dm); err != nil {
		return nil, err
	}
	return &dm, nil
}
