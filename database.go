package qb

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	// We assume sqlite
	_ "github.com/mattn/go-sqlite3"
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

// Open initializes the database
func Open(ctx context.Context, conn string) (*DB, error) {
	var err error

	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return &DB{}, err
	}

	if err = db.PingContext(ctx); err != nil {
		return &DB{}, err
	}

	return &DB{db}, nil
}

// DB represents the database
type DB struct {
	*sql.DB
}

// Delete creates and returns a new instance of DeleteQuery for the specified table
func (db *DB) Delete(ctx context.Context) *DeleteQuery {
	if tx := GetTxCtx(ctx); tx != nil {
		return tx.Delete(ctx)
	}
	return &DeleteQuery{
		ctx:    ctx,
		runner: db,
	}
}

// Insert creates and returns a new instance of InsertQuery for the specified table
func (db *DB) Insert(ctx context.Context) *InsertQuery {
	if tx := GetTxCtx(ctx); tx != nil {
		return tx.Insert(ctx)
	}
	return &InsertQuery{
		ctx:    ctx,
		runner: db,
	}
}

// Select creates and returns a new instance of SelectQuery for the specified table
func (db *DB) Select(ctx context.Context) *SelectQuery {
	if tx := GetTxCtx(ctx); tx != nil {
		return tx.Select(ctx)
	}
	return &SelectQuery{
		ctx:    ctx,
		runner: db,
	}
}

// Update creates and returns a new instance of UpdateQuery for the specified table
func (db *DB) Update(ctx context.Context) *UpdateQuery {
	if tx := GetTxCtx(ctx); tx != nil {
		return tx.Update(ctx)
	}
	return &UpdateQuery{
		ctx:    ctx,
		runner: db,
	}
}

// BeginTx starts a transaction. The default isolation level is dependent on the driver
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)

	return &Tx{tx}, err
}

// ExecContext executes a query without returning any rows. The args are for any placeholder parameters in the query
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := GetLoggerCtx(ctx)
	start := time.Now()
	result, err := db.DB.ExecContext(ctx, query, args...)
	end := time.Now()
	logger(end.Sub(start), "%s -- %v", query, args)
	return result, err
}

// QueryContext executes a query that returns rows, typically a SELECT. The args are for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	logger := GetLoggerCtx(ctx)
	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	end := time.Now()
	logger(end.Sub(start), "%s -- %v", query, args)
	return rows, err
}
