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
