package qb

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
)

// DeleteQuery represents a DELETE sql query
type DeleteQuery struct {
	runner
	ctx    context.Context
	table  string
	wheres []string
	params []interface{}
}

// From is used to set the table to delete from
func (q *DeleteQuery) From(table string) *DeleteQuery {
	q.table = table
	return q
}

// Where adds a where clause to the select query using *AND* strategy
func (q *DeleteQuery) Where(condition string, params ...interface{}) *DeleteQuery {
	q.wheres = append(q.wheres, condition)
	q.params = append(q.params, params...)
	return q
}

// Exec executes the query
func (q *DeleteQuery) Exec() (sql.Result, error) {
	return exec(q.ctx, q.runner, q)
}

// Params returns the parameters for this query
func (q *DeleteQuery) Params() []interface{} {
	return q.params
}

// Build renders the DELETE query as a string
func (q *DeleteQuery) Build(buf *bytes.Buffer) error {
	buf.WriteString("DELETE FROM ")
	buf.WriteString(q.table)

	if len(q.wheres) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(q.wheres, " AND "))
	}

	return nil
}
