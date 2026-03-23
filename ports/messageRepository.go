package ports

import "crisplite/domain"

type MessageRepository interface {
	BulkMessageInsert(batch []*domain.Message) error
}
