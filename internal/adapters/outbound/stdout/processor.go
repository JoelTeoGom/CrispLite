package database

import (
	"fmt"
	"time"

	"crisplite/internal/domain"
)

type PostgresAdapter struct{}

func NewPostgresAdapter() *PostgresAdapter {
	return &PostgresAdapter{}
}

func (p *PostgresAdapter) SaveBatch(batch []*domain.Message) error {
	fmt.Printf("Processing batch of %d messages\n", len(batch))
	for _, msg := range batch {
		fmt.Printf("Message from %s to %s: %s at %s\n", msg.SenderId, msg.ReceiverId, msg.Content, msg.Timestamp.Format(time.RFC3339))
	}
	return nil
}
