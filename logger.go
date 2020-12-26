package qb

import "context"

var (
	// LoggerContextKey store a logger in the context
	LoggerContextKey = contextKey("logger")
)

// Logger logs
type Logger func(format string, v ...interface{})

type contextKey string

func (c contextKey) String() string {
	return "qb context key " + string(c)
}

// GetLogCtx extracts a qb logger from the context
func GetLogCtx(ctx context.Context) Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(Logger); ok {
		return logger
	}

	return func(format string, v ...interface{}) {}
}
