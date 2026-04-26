package qb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// UpdateQuery represents a UPDATE sql query
type UpdateQuery struct {
	runner runner
	whereClause
	table     string
	columns   []string
	values    []interface{}
	returning []string
}

// Table is used to set the table to update
func (q *UpdateQuery) Table(table string) *UpdateQuery {
	q.table = table
	return q
}

// Set adds a column = value statement to the UPDATE query's SET clause
func (q *UpdateQuery) Set(column string, values ...interface{}) *UpdateQuery {
	q.columns = append(q.columns, column)
	q.values = append(q.values, values...)
	return q
}

// Where adds a where clause to the update query using *AND* strategy
func (q *UpdateQuery) Where(condition string, params ...interface{}) *UpdateQuery {
	q.addWhere(condition, params...)
	return q
}

// Returning specifies which columns to return after the UPDATE is successful
func (q *UpdateQuery) Returning(returning ...string) *UpdateQuery {
	q.returning = returning
	return q
}

// Exec executes the query
func (q *UpdateQuery) Exec(ctx context.Context) (sql.Result, error) {
	return exec(ctx, q.runner, q)
}

// Params returns all parameters for the query
func (q *UpdateQuery) Params() []interface{} {
	return append(q.values, q.whereClause.params...)
}

// Build renders the UPDATE query as a string
func (q *UpdateQuery) Build(buf *bytes.Buffer) error {
	buf.WriteString("UPDATE ")
	buf.WriteString(q.table)

	buf.WriteString(" SET ")
	sets := []string{}
	for _, column := range q.columns {
		sets = append(sets, fmt.Sprintf("%s = ?", column))
	}
	buf.WriteString(strings.Join(sets, ", "))

	q.writeWhere(buf)

	if len(q.returning) > 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(q.returning, ", "))
	}

	return nil
}
