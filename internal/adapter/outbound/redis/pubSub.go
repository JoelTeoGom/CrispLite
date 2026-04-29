package redis

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"

	"github.com/redis/go-redis/v9"
)

type PubSub struct {
	RedisClient *redis.Client
	PubSub      *redis.PubSub
	ChatService inbound.ChatService
	logger      outbound.Logger
}

func NewPubSub(ctx context.Context, logger outbound.Logger, redisCfg domain.RedisConfig) (*PubSub, error) {
	client, err := NewClient(ctx, logger, redisCfg)
	if err != nil {
		return nil, err
	}
	pubsub := client.Subscribe(ctx)
	return &PubSub{RedisClient: client, PubSub: pubsub, logger: logger}, nil
}

func (p *PubSub) StartPubSub(ctx context.Context, channel string) error {
	ch := p.PubSub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			go func() {
				domainMsg, err := MapRedisMessageToDomainMessage(msg)
				if err != nil {
					p.logger.Error(ctx, err)
					return
				}
				if err := p.ChatService.Deliver(ctx, domainMsg); err != nil {
					p.logger.Error(ctx, err)
				}
			}()
			p.logger.Info(ctx, "Received message from channel "+channel+": "+msg.Payload)
		}
	}
}

func (p *PubSub) Subscribe(ctx context.Context, channel string) error {
	return p.PubSub.Subscribe(ctx, channel)
}

func (p *PubSub) Publish(ctx context.Context, channel string, message any) error {
	return p.RedisClient.Publish(ctx, channel, message).Err()
}

func (p *PubSub) Close() error {
	return p.PubSub.Close()
}

func (p *PubSub) Unsubscribe(ctx context.Context, channel string) error {
	return p.PubSub.Unsubscribe(ctx, channel)
}
