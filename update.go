package qb

import (
	"bytes"
	"fmt"
	"strings"
)

// UpdateQuery represents a UPDATE sql query
type UpdateQuery struct {
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

// Params returns all parameters for the query
func (q *UpdateQuery) Params() []interface{} {
	total := len(q.values) + len(q.whereClause.params)
	if total == 0 {
		return nil
	}
	p := make([]interface{}, 0, total)
	p = append(p, q.values...)
	p = append(p, q.whereClause.params...)
	return p
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
