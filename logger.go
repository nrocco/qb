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

// GetLoggerCtx extracts a qb.Logger from the context, or nil if none is set.
func GetLoggerCtx(ctx context.Context) Logger {
	logger, _ := ctx.Value(loggerContextKey).(Logger)
	return logger
}

// WithLogger adds a qb.Logger to the context
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}
