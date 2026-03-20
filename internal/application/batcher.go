package application

import (
	"time"

	"crisplite/internal/domain"
	"crisplite/internal/ports"
)

func Batcher(messages <-chan domain.Message, batchSize int, batchInterval time.Duration, processor ports.BatchProcessor) {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	var batch []*domain.Message
	for {
		select {
		case <-ticker.C:
			if len(batch) > 0 {
				processor.ProcessBatch(batch)
				batch = nil
			}
		case msg, ok := <-messages:
			if !ok {
				if len(batch) > 0 {
					processor.ProcessBatch(batch)
				}
				return
			}
			batch = append(batch, &msg)
			if len(batch) >= batchSize {
				processor.ProcessBatch(batch)
				batch = nil
			}
		}
	}
}
