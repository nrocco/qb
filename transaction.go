package qb

import (
	"database/sql"
)

// Tx represents a transaction in a database
type Tx struct {
	*sql.Tx
	logger Logger
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

// Log uses the logger function of the database to log internals
func (tx *Tx) Log(format string, v ...interface{}) {
	tx.logger(format, v...)
}