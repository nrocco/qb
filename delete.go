package qb

import (
	"bytes"
)

// DeleteQuery represents a DELETE sql query
type DeleteQuery struct {
	whereClause
	table string
}

// From is used to set the table to delete from
func (q *DeleteQuery) From(table string) *DeleteQuery {
	q.table = table
	return q
}

// Where adds a where clause to the delete query using *AND* strategy
func (q *DeleteQuery) Where(condition string, params ...interface{}) *DeleteQuery {
	q.addWhere(condition, params...)
	return q
}

// Params returns the parameters for this query
func (q *DeleteQuery) Params() []interface{} {
	return q.whereClause.params
}

// Build renders the DELETE query as a string
func (q *DeleteQuery) Build(buf *bytes.Buffer) error {
	buf.WriteString("DELETE FROM ")
	buf.WriteString(q.table)

	q.writeWhere(buf)

	return nil
}
