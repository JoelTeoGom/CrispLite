package outbound

import (
	"context"
	"crisplite/internal/domain"
)

type MessageRepository interface {
	BulkMessageInsert(ctx context.Context, batch []*domain.Message) error
}
