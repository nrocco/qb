package qb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
)

// SelectQuery represents a SELECT sql query
type SelectQuery struct {
	runner
	ctx      context.Context
	table    string
	wheres   []string
	params   []interface{}
	columns  []string
	joins    []string
	limit    string
	offset   string
	cte      string
	orderBys []string
	groupBys []string
}

// From is used to set the table to select from
func (q *SelectQuery) From(table string) *SelectQuery {
	q.table = table
	return q
}

// Columns determines with columns to select
func (q *SelectQuery) Columns(columns ...string) *SelectQuery {
	q.columns = columns
	return q
}

// Join adds a join to the select query
func (q *SelectQuery) Join(join string, params ...interface{}) *SelectQuery {
	q.joins = append(q.joins, join)
	q.params = append(q.params, params...)
	return q
}

// Where adds a where clause to the select query using *AND* strategy
func (q *SelectQuery) Where(condition string, params ...interface{}) *SelectQuery {
	q.wheres = append(q.wheres, condition)
	q.params = append(q.params, params...)
	return q
}

// OrderBy adds an ORDER BY clause to the SELECT query
func (q *SelectQuery) OrderBy(column string, direction string) *SelectQuery {
	q.orderBys = append(q.orderBys, column+" "+direction)
	return q
}

// GroupBy adds a GROUP BY clause to the SELECT query
func (q *SelectQuery) GroupBy(column string) *SelectQuery {
	q.groupBys = append(q.groupBys, column)
	return q
}

// Limit adds a LIMIT clause to the SELECT query
func (q *SelectQuery) Limit(limit int) *SelectQuery {
	q.limit = fmt.Sprintf("%d", limit)
	return q
}

// Offset adds a OFFSET clause to the SELECT query
func (q *SelectQuery) Offset(offset int) *SelectQuery {
	q.offset = fmt.Sprintf("%d", offset)
	return q
}

// With adds a Common Table Expressions to the beginning of the query
func (q *SelectQuery) With(cte string, params ...interface{}) *SelectQuery {
	q.cte = cte
	q.params = append(q.params, params...)
}

// Load will execute the query and scan the result into the given struct
func (q *SelectQuery) Load(value interface{}) (int, error) {
	return query(q.ctx, q.runner, q, value)
}

// LoadValue will execute the query and scan the scalar result into the given variable
func (q *SelectQuery) LoadValue(value interface{}) error {
	rows, err := query(q.ctx, q.runner, q, value)
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("no records returned")
	}

	return nil
}

// Params returns the parameters for this query
func (q *SelectQuery) Params() []interface{} {
	return q.params
}

// Build renders the SELECT query as a string
func (q *SelectQuery) Build(buf *bytes.Buffer) error {
	if q.cte {
		buf.WriteString(q.cte)
		buf.WriteString( " ")
	}

	buf.WriteString("SELECT ")

	if len(q.columns) > 0 {
		buf.WriteString(strings.Join(q.columns, ", "))
	} else {
		buf.WriteString("*")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(q.table)

	for _, join := range q.joins {
		buf.WriteString(" ")
		buf.WriteString(join)
	}

	if len(q.wheres) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(q.wheres, " AND "))
	}

	if len(q.groupBys) != 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(q.groupBys, ", "))
	}

	if len(q.orderBys) != 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(q.orderBys, ", "))
	}

	if q.limit != "" {
		buf.WriteString(" LIMIT ")
		buf.WriteString(q.limit)
	}

	if q.offset != "" {
		buf.WriteString(" OFFSET ")
		buf.WriteString(q.offset)
	}

	return nil
}
