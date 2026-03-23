// internal/adapter/postgres/pool.go
package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			poolErr = fmt.Errorf("parsing pg config: %w", err)
			return
		}

		config.MaxConns = 20
		config.MinConns = 2
		config.MaxConnLifetime = 30 * time.Minute
		config.MaxConnIdleTime = 5 * time.Minute

		p, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			poolErr = fmt.Errorf("creating pg pool: %w", err)
			return
		}

		if err := p.Ping(ctx); err != nil {
			p.Close()
			poolErr = fmt.Errorf("pinging pg: %w", err)
			return
		}

		pool = p
	})

	return pool, poolErr
}

func GetPool() *pgxpool.Pool {
	return pool
}
