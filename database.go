package qb

import (
	"bytes"
	"context"
	"database/sql"

	// We assume sqlite
	_ "modernc.org/sqlite"
)

type runner interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Builder is the interface for all query builders
type Builder interface {
	Build(buf *bytes.Buffer) error
	Params() []interface{}
}

type contextKey string

func (c contextKey) String() string {
	return "qb context key " + string(c)
}

// Executor holds a runner and provides factory methods for all query types.
type Executor struct {
	runner runner
}

// Select creates and returns a new SelectQuery
func (e *Executor) Select() *SelectQuery {
	return &SelectQuery{runner: e.runner}
}

// Insert creates and returns a new InsertQuery
func (e *Executor) Insert() *InsertQuery {
	return &InsertQuery{runner: e.runner}
}

// Update creates and returns a new UpdateQuery
func (e *Executor) Update() *UpdateQuery {
	return &UpdateQuery{runner: e.runner}
}

// Delete creates and returns a new DeleteQuery
func (e *Executor) Delete() *DeleteQuery {
	return &DeleteQuery{runner: e.runner}
}

// Open initializes the database
func Open(ctx context.Context, conn string) (*DB, error) {
	db, err := sql.Open("sqlite", conn)
	if err != nil {
		return &DB{}, err
	}

	if err = db.PingContext(ctx); err != nil {
		return &DB{}, err
	}

	return &DB{DB: db, Executor: Executor{runner: db}}, nil
}

// DB represents the database
type DB struct {
	*sql.DB
	Executor
}

// For returns the Executor for the active transaction in ctx, or the DB's own Executor.
func (db *DB) For(ctx context.Context) *Executor {
	if tx := GetTxCtx(ctx); tx != nil {
		return &tx.Executor
	}
	return &db.Executor
}

// BeginTx starts a transaction. The default isolation level is dependent on the driver.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	return &Tx{Tx: tx, Executor: Executor{runner: tx}}, err
}
