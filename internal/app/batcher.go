package app

import (
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
}

func NewBatcher(msgs <-chan domain.Message, size int, interval time.Duration, proc outbound.MessageRepository) *Batcher {
	return &Batcher{
		batchSize: size,
		interval:  interval,
		processor: proc,
		messages:  msgs,
		done:      make(chan struct{}),
	}
}

func (b *Batcher) Start() {
	go b.run()
}

func (b *Batcher) Stop() {
	defer close(b.done)
}

func (b *Batcher) run() {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	var batch []*domain.Message
	for {
		select {
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
