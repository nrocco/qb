package qb

import (
	"bytes"
	"context"
	"database/sql"
	"reflect"
	"strings"
)

// InsertQuery represents a INSERT sql query
type InsertQuery struct {
	runner
	ctx            context.Context
	orIgnore       bool
	table          string
	columns        []string
	values         []interface{}
	conflictColumn string
	conflictSets   string
	returning      []string
	recordID       reflect.Value
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

// Values determines the columns to insert
func (q *InsertQuery) Values(values ...interface{}) *InsertQuery {
	q.values = values
	return q
}

// OnConflict specifies what to do if there is conflict
func (q *InsertQuery) OnConflict(column string, sets string) *InsertQuery {
	q.conflictColumn = column
	q.conflictSets = sets
	return q
}

// Returning specifies with columns to return after the INSERT is successful
func (q *InsertQuery) Returning(returning ...string) *InsertQuery {
	q.returning = returning
	return q
}

// Exec executes the query
func (q *InsertQuery) Exec() (sql.Result, error) {
	result, err := exec(q.ctx, q.runner, q)
	if err != nil {
		return nil, err
	}

	if q.recordID.IsValid() {
		if id, err := result.LastInsertId(); err == nil {
			q.recordID.SetInt(id)
		}
	}

	return result, nil
}

// Record scans the result of the query into the given struct
func (q *InsertQuery) Record(structValue interface{}) {
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

		for _, name := range []string{"Id", "ID"} {
			field := value.FieldByName(name)
			if field.IsValid() && field.Kind() == reflect.Int64 {
				q.recordID = field
				break
			}
		}

		q.Values(values...)
	}
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
