package postgres

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	pool   *pgxpool.Pool
	logger outbound.Logger
}

func NewMessageRepo(pool *pgxpool.Pool, logger outbound.Logger) *MessageRepo {
	return &MessageRepo{pool: pool, logger: logger}
}

func (p *MessageRepo) BulkMessageInsert(ctx context.Context, batch []*domain.Message) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	columns := 5
	placeholders := make([]string, len(batch))
	args := make([]any, 0, len(batch)*columns)
	for i, m := range batch {
		base := i * columns
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5)
		args = append(args, m.ConversationId, m.SenderId, m.ReceiverId, m.Timestamp, m.Content)
	}

	query := "INSERT INTO messages (conversation_id, sender_id, receiver_id, timestamp, content) VALUES " +
		strings.Join(placeholders, ", ")

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		p.logger.Error(ctx, err)
		return err
	}

	for _, message := range batch {
		_, err := tx.Exec(ctx, "UPDATE conversations SET last_message_at = $1 WHERE id = $2 AND (last_message_at IS NULL OR last_message_at < $1)", message.Timestamp, message.ConversationId)
		if err != nil {
			p.logger.Error(ctx, err)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *MessageRepo) UpdateConversationLastMessage(ctx context.Context, message domain.Message) error {
	_, err := p.pool.Exec(ctx, "UPDATE conversations SET last_message_at = $1 WHERE id = $2 AND (last_message_at IS NULL OR last_message_at < $1)", message.Timestamp, message.ConversationId)
	if err != nil {
		p.logger.Error(ctx, err)
		return err
	}
	return nil
}
