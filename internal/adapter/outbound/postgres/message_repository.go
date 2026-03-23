package postgres

import (
	"crisplite/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	pool *pgxpool.Pool
}

func NewMessageRepo(pool *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{pool: pool}
}

func (p *MessageRepo) BulkMessageInsert(batch []*domain.Message) error {
	return nil
}
