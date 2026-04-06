package locallogger

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
)

type LocalLogger struct{}

func NewLocalLogger() *LocalLogger {
	return &LocalLogger{}
}

func (l *LocalLogger) Info(ctx context.Context, msg string) {
	log.Printf("[INFO] %s", msg)
}

func (l *LocalLogger) Warning(ctx context.Context, msg string) {
	log.Printf("[WARNING] %s", msg)
}

func (l *LocalLogger) Error(ctx context.Context, err error) {
	log.Printf("[ERROR] %v", err)
	debug.PrintStack()
}

func (l *LocalLogger) ErrorWithVar(ctx context.Context, err error, data map[string]any) {
	log.Printf("[ERROR] %v %s", err, formatData(data))
	debug.PrintStack()
}

func (l *LocalLogger) Debug(ctx context.Context, msg string) {
	log.Printf("[DEBUG] %s", msg)
}

func (l *LocalLogger) DebugWithVar(ctx context.Context, msg string, data map[string]any) {
	log.Printf("[DEBUG] %s %s", msg, formatData(data))
}

func formatData(data map[string]any) string {
	if len(data) == 0 {
		return ""
	}
	result := ""
	for k, v := range data {
		result += fmt.Sprintf("%s=%v ", k, v)
	}
	return result
}
