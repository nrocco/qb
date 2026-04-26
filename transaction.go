package qb

import (
	"context"
	"database/sql"
)

var (
	txContextKey = contextKey("tx")
)

// GetTxCtx extracts a qb.Tx from the context
func GetTxCtx(ctx context.Context) *Tx {
	tx, _ := ctx.Value(txContextKey).(*Tx)
	return tx
}

// WithTx adds a qb.Tx to the context
func WithTx(ctx context.Context, tx *Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

// Tx represents a transaction in a database
type Tx struct {
	*sql.Tx
	Executor
}
