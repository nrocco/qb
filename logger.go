package qb

import (
	"context"
	"time"
)

var (
	loggerContextKey = contextKey("logger")
)

// Logger logs
type Logger func(ctx context.Context, duration time.Duration, format string, v ...interface{})

// GetLoggerCtx extracts a qb.Logger from the context
func GetLoggerCtx(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return logger
	}

	return func(ctx context.Context, duration time.Duration, format string, v ...interface{}) {}
}

// WitLogger adds a qb.Logger to the context
func WitLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}
