package outbound

import "context"

type PubSub interface {
	Start(ctx context.Context, channel string) error
	Subscribe(ctx context.Context, channel string) error
	Publish(ctx context.Context, channel string, message any) error
	Unsubscribe(ctx context.Context, channel string) error
	Close() error
}
