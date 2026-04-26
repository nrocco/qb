package qb

import (
	"context"
	"database/sql"
	"errors"
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
}

// Select creates and returns a new SelectQuery
func (tx *Tx) Select() *SelectQuery { return &SelectQuery{} }

// Insert creates and returns a new InsertQuery
func (tx *Tx) Insert() *InsertQuery { return &InsertQuery{} }

// Update creates and returns a new UpdateQuery
func (tx *Tx) Update() *UpdateQuery { return &UpdateQuery{} }

// Delete creates and returns a new DeleteQuery
func (tx *Tx) Delete() *DeleteQuery { return &DeleteQuery{} }

// Exec executes a write query within the transaction
func (tx *Tx) Exec(ctx context.Context, b Builder) (sql.Result, error) {
	return exec(ctx, tx.Tx, b)
}

// Load executes a read query within the transaction and scans the results into dest
func (tx *Tx) Load(ctx context.Context, b Builder, dest interface{}) (int, error) {
	return query(ctx, tx.Tx, b, dest)
}

// LoadValue executes a read query within the transaction and scans the scalar result into dest
func (tx *Tx) LoadValue(ctx context.Context, b Builder, dest interface{}) error {
	rows, err := tx.Load(ctx, b, dest)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no records returned")
	}
	return nil
}
