package main

import (
	"crisplite/database"
	"crisplite/domain"
	"crisplite/ports"
	"fmt"
	"net/http"
	"time"
)

type Batcher struct {
	batchSize int
	interval  time.Duration
	processor ports.MessageRepository
	messages  <-chan domain.Message
	done      chan struct{}
}

func NewBatcher(msgs <-chan domain.Message, size int, interval time.Duration, proc ports.MessageRepository) *Batcher {
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

func main() {
	batchSize := 100
	msgChannel := make(chan domain.Message, batchSize)
	defer close(msgChannel)

	postgres := database.NewPostgresAdapter()

	batcher := NewBatcher(msgChannel, batchSize, 200*time.Millisecond, postgres)
	//http.HandleFunc("/ws/v1/chat", wsHandler.ChatHandler)

	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
