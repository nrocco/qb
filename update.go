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
	runner
	table     string
	wheres    []string
	params    []interface{}
	columns   []string
	values    []interface{}
	returning []string
}

// Set adds a column = value statement to the UPDATE querie's SET clause
func (q *UpdateQuery) Set(column string, values ...interface{}) *UpdateQuery {
	q.columns = append(q.columns, column)
	q.values = append(q.values, values...)
	return q
}

// Where adds a where clause to the update query using *AND* strategy
func (q *UpdateQuery) Where(condition string, params ...interface{}) *UpdateQuery {
	q.wheres = append(q.wheres, condition)
	q.params = append(q.params, params...)
	return q
}

// Returning specifies with columns to return after the UPDATE is successful
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
	return append(q.values, q.params...)
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

	if len(q.wheres) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(q.wheres, " AND "))
	}

	if len(q.returning) > 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(q.returning, ", "))
	}

	return nil
}
