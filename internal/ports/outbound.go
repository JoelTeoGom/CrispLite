package ports

import "crisplite/internal/domain"

type BatchProcessor interface {
	ProcessBatch(batchMessages []*domain.Message) error
}
