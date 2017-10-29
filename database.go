package qb

import (
	"bytes"
	"database/sql"

	// We assume sqlite
	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database
type DB struct {
	*sql.DB
}

type runner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Builder is implemented by all types of query builders
type Builder interface {
	Build(buf *bytes.Buffer) error
	Params() []interface{}
}

// Open initializes the database
func Open(conn string) (*DB, error) {
	var err error

	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return &DB{}, err
	}

	if err = db.Ping(); err != nil {
		return &DB{}, err
	}

	return &DB{db}, nil
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
