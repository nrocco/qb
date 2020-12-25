package qb

import (
	"context"
	"database/sql"
	"fmt"
)

// Tx represents a transaction in a database
type Tx struct {
	*sql.Tx
}

// Commit commits the transaction
func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
}

// Rollback aborts the transaction
func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

// Delete creates and returns a new instance of DeleteQuery for the specified table
func (tx *Tx) Delete(table string) *DeleteQuery {
	return &DeleteQuery{
		runner: tx,
		table:  table,
	}
}

// Insert creates and returns a new instance of InsertQuery for the specified table
func (tx *Tx) Insert(table string) *InsertQuery {
	return &InsertQuery{
		runner: tx,
		table:  table,
	}
}

// Select creates and returns a new instance of SelectQuery for the specified table
func (tx *Tx) Select(table string) *SelectQuery {
	return &SelectQuery{
		runner: tx,
		table:  table,
	}
}

// Update creates and returns a new instance of UpdateQuery for the specified table
func (tx *Tx) Update(table string) *UpdateQuery {
	return &UpdateQuery{
		runner: tx,
		table:  table,
	}
}

// Savepoint starts a savepoint. TODO remove this???
func (tx *Tx) Savepoint(ctx context.Context, name string) error {
	_, err := tx.ExecContext(ctx, fmt.Sprintf("SAVEPOINT %s", name))

	return err
}

// ReleaseSavepoint commits a savepoint. TODO remove this???
func (tx *Tx) ReleaseSavepoint(ctx context.Context, name string) error {
	_, err := tx.ExecContext(ctx, fmt.Sprintf("RELEASE SAVEPOINT %s", name))

	return err
}

// RollbackSavepoint rolls back a savepoint. TODO remove this???
func (tx *Tx) RollbackSavepoint(ctx context.Context, name string) error {
	_, err := tx.ExecContext(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name))

	return err
}
