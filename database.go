package qb

import (
	"bytes"
	"context"
	"database/sql"

	// We assume sqlite
	_ "github.com/mattn/go-sqlite3"
)

// Logger logs
type Logger func(format string, v ...interface{})

// DB represents the database
type DB struct {
	*sql.DB
	logger Logger
}

type runner interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Builder is the interface for all query builders
type Builder interface {
	Build(buf *bytes.Buffer) error
	Params() []interface{}
}

// Open initializes the database
func Open(ctx context.Context, conn string, logger Logger) (*DB, error) {
	var err error

	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return &DB{}, err
	}

	if err = db.PingContext(ctx); err != nil {
		return &DB{}, err
	}

	return &DB{db, logger}, nil
}

// ExecContext executes the given SQL query against the database
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db.logger != nil {
		db.logger("%s -- %v", query, args)
	}

	return db.DB.ExecContext(ctx, query, args...)
}

// QueryContext executes the given SQL query against the database
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.logger != nil {
		db.logger("%s -- %v", query, args)
	}

	return db.DB.QueryContext(ctx, query, args...)
}

// Delete creates and returns a new instance of DeleteQuery for the specified table
func (db *DB) Delete(table string) *DeleteQuery {
	return &DeleteQuery{
		runner: db,
		table:  table,
	}
}

// Insert creates and returns a new instance of InsertQuery for the specified table
func (db *DB) Insert(table string) *InsertQuery {
	return &InsertQuery{
		runner: db,
		table:  table,
	}
}

// Select creates and returns a new instance of SelectQuery for the specified table
func (db *DB) Select(table string) *SelectQuery {
	return &SelectQuery{
		runner: db,
		table:  table,
	}
}

// Update creates and returns a new instance of UpdateQuery for the specified table
func (db *DB) Update(table string) *UpdateQuery {
	return &UpdateQuery{
		runner: db,
		table:  table,
	}
}

// BeginTx starts a transaction. The default isolation level is dependent on the driver
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)

	return &Tx{tx}, err
}
