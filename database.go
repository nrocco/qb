package qb

import (
	"bytes"
	"database/sql"
	"fmt"

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
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Builder is the interface for all query builders
type Builder interface {
	Build(buf *bytes.Buffer) error
	Params() []interface{}
}

// Open initializes the database
func Open(conn string, logger Logger) (*DB, error) {
	var err error

	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return &DB{}, err
	}

	if err = db.Ping(); err != nil {
		return &DB{}, err
	}

	return &DB{db, logger}, nil
}

// Exec executes the given SQL query against the database
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	db.logger("%s -- %v", query, args)

	return db.DB.Exec(query, args...)
}

// Query executes the given SQL query against the database
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db.logger("%s -- %v", query, args)

	return db.DB.Query(query, args...)
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

// Savepoint starts a savepoint.
func (db *DB) Savepoint(name string) error {
	_, err := db.Exec(fmt.Sprintf("SAVEPOINT %s", name))

	return err
}

// ReleaseSavepoint commits a savepoint.
func (db *DB) ReleaseSavepoint(name string) error {
	_, err := db.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", name))

	return err
}

// RollbackSavepoint rolls back a savepoint.
func (db *DB) RollbackSavepoint(name string) error {
	_, err := db.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name))

	return err
}

// Begin starts a transaction. The default isolation level is dependent on the driver
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()

	return &Tx{tx}, err
}
