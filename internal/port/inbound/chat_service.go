package inbound

import (
	"context"
	"crisplite/internal/domain"
)

type ChatService interface {
	Send(ctx context.Context, msg *domain.Message) error
	Deliver(ctx context.Context, msg *domain.Message) error
}
