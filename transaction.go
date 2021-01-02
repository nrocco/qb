package qb

import (
	"context"
	"database/sql"
	"time"
)

var (
	txContextKey = contextKey("tx")
)

// GetTxCtx extracts a qb.Tx from the context
func GetTxCtx(ctx context.Context) *Tx {
	tx, _ := ctx.Value(txContextKey).(*Tx)
	return tx
}

// WitTx adds a qb.Tx to the context
func WitTx(ctx context.Context, tx *Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

// Tx represents a transaction in a database
type Tx struct {
	*sql.Tx
}

// Delete creates and returns a new instance of DeleteQuery for the specified table
func (tx *Tx) Delete(ctx context.Context) *DeleteQuery {
	return &DeleteQuery{
		ctx:    ctx,
		runner: tx,
	}
}

// Insert creates and returns a new instance of InsertQuery for the specified table
func (tx *Tx) Insert(ctx context.Context) *InsertQuery {
	return &InsertQuery{
		ctx:    ctx,
		runner: tx,
	}
}

// Select creates and returns a new instance of SelectQuery for the specified table
func (tx *Tx) Select(ctx context.Context) *SelectQuery {
	return &SelectQuery{
		ctx:    ctx,
		runner: tx,
	}
}

// Update creates and returns a new instance of UpdateQuery for the specified table
func (tx *Tx) Update(ctx context.Context) *UpdateQuery {
	return &UpdateQuery{
		ctx:    ctx,
		runner: tx,
	}
}

// ExecContext executes a query without returning any rows. The args are for any placeholder parameters in the query
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := GetLoggerCtx(ctx)
	start := time.Now()
	result, err := tx.Tx.ExecContext(ctx, query, args...)
	end := time.Now()
	logger(end.Sub(start), "%s -- %v", query, args)
	return result, err
}

// QueryContext executes a query that returns rows, typically a SELECT. The args are for any placeholder parameters in the query.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	logger := GetLoggerCtx(ctx)
	start := time.Now()
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	end := time.Now()
	logger(end.Sub(start), "%s -- %v", query, args)
	return rows, err
}
