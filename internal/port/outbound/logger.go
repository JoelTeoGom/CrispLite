package outbound

import "context"

type Logger interface {
	Info(ctx context.Context, msg string)
	Warning(ctx context.Context, msg string)
	Error(ctx context.Context, err error)
	ErrorWithVar(ctx context.Context, err error, data map[string]any)
	Debug(ctx context.Context, msg string)
	DebugWithVar(ctx context.Context, msg string, data map[string]any)
}
