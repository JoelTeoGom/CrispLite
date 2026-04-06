package postgres

import (
	"context"
	"crisplite/internal/port/outbound"
	"fmt"

	"github.com/jackc/pgx/v5/tracelog"
)

type PgxLogAdapter struct {
	logger outbound.Logger
}

func NewPgxLogAdapter(logger outbound.Logger) *PgxLogAdapter {
	return &PgxLogAdapter{logger: logger}
}

func (a *PgxLogAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case tracelog.LogLevelError:
		a.logger.ErrorWithVar(ctx, fmt.Errorf(msg), data)
	case tracelog.LogLevelWarn:
		a.logger.Warning(ctx, msg)
	case tracelog.LogLevelInfo:
		a.logger.Info(ctx, msg)
	case tracelog.LogLevelDebug:
		a.logger.DebugWithVar(ctx, msg, data)
	}
}
