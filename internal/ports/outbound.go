package ports

import "crisplite/internal/domain"

type MessageRepository interface {
	SaveBatch(batchMessages []*domain.Message) error
}
