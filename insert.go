package qb

import (
	"bytes"
	"database/sql"
	"reflect"
	"strings"
)

// InsertQuery represents a INSERT sql query
type InsertQuery struct {
	runner
	table     string
	columns   []string
	values    []interface{}
	returning []string
	recordID  reflect.Value
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

// Returning can be used to choose which columns to return after the INSERT is succesful
func (q *InsertQuery) Returning(returning ...string) *InsertQuery {
	q.returning = returning
	return q
}

func (q *InsertQuery) Exec() (sql.Result, error) {
	result, err := exec(q.runner, q)
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

func (q *InsertQuery) Record(structValue interface{}) {
	v := reflect.Indirect(reflect.ValueOf(structValue))

	if v.Kind() == reflect.Struct {
		var value []interface{}
		m := structMap(v)
		for _, key := range q.columns {
			if val, ok := m[key]; ok {
				value = append(value, val.Interface())
			} else {
				value = append(value, nil)
			}
		}

		for _, name := range []string{"Id", "ID"} {
			field := v.FieldByName(name)
			if field.IsValid() && field.Kind() == reflect.Int64 {
				q.recordID = field
				break
			}
		}

		q.Values(value...)
	}
}

// Build renders the INSERT query as a string
func (q *InsertQuery) Build(buf *bytes.Buffer) error {
	buf.WriteString("INSERT INTO ")
	buf.WriteString(q.table)
	if len(q.columns) > 0 {
		buf.WriteString(" (")
		buf.WriteString(strings.Join(q.columns, ", "))
		buf.WriteString(")")
	}

	buf.WriteString(" VALUES (")
	fuus := []string{} // TODO this can be done better
	for _, _ = range q.values {
		fuus = append(fuus, "?")
	}
	buf.WriteString(strings.Join(fuus, ", "))
	buf.WriteString(")")

	if len(q.returning) > 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(q.returning, ", "))
	}

	return nil
}

func (q *InsertQuery) Params() []interface{} {
	return q.values
}
