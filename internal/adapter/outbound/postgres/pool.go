package postgres

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

func NewPool(ctx context.Context, dbCfg domain.DatabaseConfig, logger outbound.Logger) (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		config, err := pgxpool.ParseConfig(dbCfg.URL)
		if err != nil {
			poolErr = fmt.Errorf("parsing pg config: %w", err)
			return
		}
		pgxLogger := NewPgxLogAdapter(logger)
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxLogger,
			LogLevel: tracelog.LogLevelInfo,
		}
		config.MaxConns = int32(dbCfg.MaxConns)
		config.MinConns = int32(dbCfg.MinConns)
		config.MaxConnLifetime = dbCfg.MaxConnLife
		config.MaxConnIdleTime = dbCfg.MaxConnIdle

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
