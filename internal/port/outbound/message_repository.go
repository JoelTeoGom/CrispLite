package outbound

import "crisplite/internal/domain"

type MessageRepository interface {
	BulkMessageInsert(batch []*domain.Message) error
}
