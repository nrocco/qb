package qb

import (
	"bytes"
	"reflect"
	"strings"
)

// InsertQuery represents a INSERT sql query
type InsertQuery struct {
	orIgnore       bool
	table          string
	columns        []string
	values         []interface{}
	conflictColumn string
	conflictSets   string
	returning      []string
}

// OrIgnore make the query behave using INSERT OR IGNORE INTO
func (q *InsertQuery) OrIgnore() *InsertQuery {
	q.orIgnore = true
	return q
}

// InTo is used to set the table to insert into
func (q *InsertQuery) InTo(table string) *InsertQuery {
	q.table = table
	return q
}

// Columns determines the columns to insert
func (q *InsertQuery) Columns(columns ...string) *InsertQuery {
	q.columns = columns
	return q
}

// Values determines the values to insert
func (q *InsertQuery) Values(values ...interface{}) *InsertQuery {
	q.values = values
	return q
}

// OnConflict specifies what to do if there is a conflict
func (q *InsertQuery) OnConflict(column string, sets string) *InsertQuery {
	q.conflictColumn = column
	q.conflictSets = sets
	return q
}

// Returning specifies which columns to return after the INSERT is successful
func (q *InsertQuery) Returning(returning ...string) *InsertQuery {
	q.returning = returning
	return q
}

// Record populates Values from the struct fields matching Columns
func (q *InsertQuery) Record(structValue interface{}) *InsertQuery {
	value := reflect.Indirect(reflect.ValueOf(structValue))

	if value.Kind() == reflect.Struct {
		var values []interface{}
		structFields := structMap(value)
		for _, column := range q.columns {
			if val, ok := structFields[column]; ok {
				values = append(values, val.Interface())
			} else {
				values = append(values, nil)
			}
		}
		q.Values(values...)
	}

	return q
}

// Build renders the INSERT query as a string
func (q *InsertQuery) Build(buf *bytes.Buffer) error {
	if q.orIgnore {
		buf.WriteString("INSERT OR IGNORE INTO ")
	} else {
		buf.WriteString("INSERT INTO ")
	}
	buf.WriteString(q.table)
	if len(q.columns) > 0 {
		buf.WriteString(" (")
		buf.WriteString(strings.Join(q.columns, ", "))
		buf.WriteString(")")
	}

	buf.WriteString(" VALUES (")
	fuus := []string{} // TODO make this better
	for range q.values {
		fuus = append(fuus, "?")
	}
	buf.WriteString(strings.Join(fuus, ", "))
	buf.WriteString(")")

	if len(q.returning) > 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(q.returning, ", "))
	}

	if q.conflictColumn != "" {
		buf.WriteString(" ON CONFLICT (")
		buf.WriteString(q.conflictColumn)
		buf.WriteString(") DO UPDATE SET ")
		buf.WriteString(q.conflictSets)
	}

	return nil
}

// Params returns all parameters for the query
func (q *InsertQuery) Params() []interface{} {
	return q.values
}
