package postgres

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	pool   *pgxpool.Pool
	logger outbound.Logger
}

func NewMessageRepo(pool *pgxpool.Pool, logger outbound.Logger) *MessageRepo {
	return &MessageRepo{pool: pool, logger: logger}
}

func (p *MessageRepo) BulkMessageInsert(batch []*domain.Message) error {
	return nil
}
