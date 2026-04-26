package qb

import (
	"bytes"
	"context"
	"database/sql"
	"errors"

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

// DB represents the database
type DB struct {
	*sql.DB
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

	return &DB{db}, nil
}

// BeginTx starts a transaction. The default isolation level is dependent on the driver.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	return &Tx{tx}, err
}

// Select creates and returns a new SelectQuery
func (db *DB) Select() *SelectQuery { return &SelectQuery{} }

// Insert creates and returns a new InsertQuery
func (db *DB) Insert() *InsertQuery { return &InsertQuery{} }

// Update creates and returns a new UpdateQuery
func (db *DB) Update() *UpdateQuery { return &UpdateQuery{} }

// Delete creates and returns a new DeleteQuery
func (db *DB) Delete() *DeleteQuery { return &DeleteQuery{} }

func (db *DB) runnerFor(ctx context.Context) runner {
	if tx := GetTxCtx(ctx); tx != nil {
		return tx.Tx
	}
	return db.DB
}

// Exec executes a write query, using the transaction in ctx if present
func (db *DB) Exec(ctx context.Context, b Builder) (sql.Result, error) {
	return exec(ctx, db.runnerFor(ctx), b)
}

// Load executes a read query and scans the results into dest
func (db *DB) Load(ctx context.Context, b Builder, dest interface{}) (int, error) {
	return query(ctx, db.runnerFor(ctx), b, dest)
}

// LoadValue executes a read query and scans the scalar result into dest
func (db *DB) LoadValue(ctx context.Context, b Builder, dest interface{}) error {
	rows, err := db.Load(ctx, b, dest)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no records returned")
	}
	return nil
}
