package qb

import (
	"bytes"
	"database/sql"
	"strings"
)

// DeleteQuery represents a DELETE sql query
type DeleteQuery struct {
	runner
	table  string
	wheres []string
	params []interface{}
}

// Where adds a where clause to the select query using *AND* strategy
func (q *DeleteQuery) Where(condition string, params ...interface{}) {
	q.wheres = append(q.wheres, condition)
	q.params = append(q.params, params...)
}

func (q *DeleteQuery) Exec() (sql.Result, error) {
	return exec(q.runner, q)
}

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
