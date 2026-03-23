package database

import "crisplite/domain"

type PostgresAdapter struct {
}

func NewPostgresAdapter() *PostgresAdapter {
	return &PostgresAdapter{}
}

func (p *PostgresAdapter) BulkMessageInsert(batch []*domain.Message) error
