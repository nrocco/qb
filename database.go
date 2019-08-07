package qb

import (
	"bytes"
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
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Log(format string, v ...interface{})
}

// Builder is implemented by all types of query builders
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

// Begin starts a transaction. The default isolation level is dependent on the driver
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()

	return &Tx{tx, db.logger}, err
}

// Log uses the logger function of the database to log internals
func (db *DB) Log(format string, v ...interface{}) {
	db.logger(format, v...)
}
