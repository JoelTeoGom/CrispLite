package app

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"time"
)

type Batcher struct {
	batchSize int
	interval  time.Duration
	processor outbound.MessageRepository
	messages  <-chan domain.Message
	done      chan struct{}
	logger    outbound.Logger
}

func NewBatcher(msgs <-chan domain.Message, size int, interval time.Duration, proc outbound.MessageRepository, logger outbound.Logger) *Batcher {
	return &Batcher{
		batchSize: size,
		interval:  interval,
		processor: proc,
		messages:  msgs,
		done:      make(chan struct{}),
		logger:    logger,
	}
}

func (b *Batcher) Start(ctx context.Context) {
	go b.run(ctx)
}

func (b *Batcher) Stop() {
	defer close(b.done)
}

func (b *Batcher) run(ctx context.Context) {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	var batch []*domain.Message
	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				b.processor.BulkMessageInsert(batch)
			}
		case <-b.done:
			if len(batch) > 0 {
				b.processor.BulkMessageInsert(batch)
			}
			return
		case <-ticker.C:
			if len(batch) > 0 {
				b.processor.BulkMessageInsert(batch)
				batch = nil
			}
		case msg, ok := <-b.messages:
			if !ok {
				if len(batch) > 0 {
					b.processor.BulkMessageInsert(batch)
				}
				return
			}
			batch = append(batch, &msg)
			if len(batch) >= b.batchSize {
				b.processor.BulkMessageInsert(batch)
				batch = nil
			}
		}
	}
}
