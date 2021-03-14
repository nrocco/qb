package qb

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSelectQuery(t *testing.T) {
	type test struct {
		name   string
		query  func() *SelectQuery
		result string
		values []interface{}
		err    error
	}

	var testResults = []test{
		{
			name: "select everything",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				return query
			},
			result: "SELECT * FROM fuu",
		},
		{
			name: "select specific columns",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.Columns("bar", "baz")
				return query
			},
			result: "SELECT bar, baz FROM fuu",
		},
		{
			name: "select with join",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.Join("LEFT JOIN bar ON bar.id = fuu.bar_id")
				return query
			},
			result: "SELECT * FROM fuu LEFT JOIN bar ON bar.id = fuu.bar_id",
		},
		{
			name: "select single where condition",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.Where("column1 = ?", 123)
				return query
			},
			result: "SELECT * FROM fuu WHERE column1 = ?",
			values: []interface{}{123},
		},
		{
			name: "select multiple where condition",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.Where("column1 = ?", "fuu")
				query.Where("column2 IS NULL")
				return query
			},
			result: "SELECT * FROM fuu WHERE column1 = ? AND column2 IS NULL",
			values: []interface{}{"fuu"},
		},
		{
			name: "select group by",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.GroupBy("column1")
				return query
			},
			result: "SELECT * FROM fuu GROUP BY column1",
		},
		{
			name: "select multiple group by",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.GroupBy("column1")
				query.GroupBy("column2")
				return query
			},
			result: "SELECT * FROM fuu GROUP BY column1, column2",
		},
		{
			name: "select with ordering ascending",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.OrderBy("column1", "ASC")
				return query
			},
			result: "SELECT * FROM fuu ORDER BY column1 ASC",
		},
		{
			name: "select with ordering descending",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.OrderBy("column1", "DESC")
				return query
			},
			result: "SELECT * FROM fuu ORDER BY column1 DESC",
		},
		{
			name: "select with multiple order bys",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.OrderBy("column1", "DESC")
				query.OrderBy("column2", "ASC")
				return query
			},
			result: "SELECT * FROM fuu ORDER BY column1 DESC, column2 ASC",
		},
		{
			name: "select with limit and offset",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu"}
				query.Limit(10)
				query.Offset(10)
				return query
			},
			result: "SELECT * FROM fuu LIMIT 10 OFFSET 10",
		},
		{
			name: "select with many things combined",
			query: func() *SelectQuery {
				query := &SelectQuery{table: "fuu f"}
				query.Columns("COUNT(f.id)")
				query.Where("f.name = ?", "something")
				query.OrderBy("f.name", "ASC")
				query.Join("LEFT JOIN bar b ON b.fuu_id = f.id")
				query.GroupBy("b.id")
				query.Limit(5)
				query.Offset(0)
				return query
			},
			result: "SELECT COUNT(f.id) FROM fuu f LEFT JOIN bar b ON b.fuu_id = f.id WHERE f.name = ? GROUP BY b.id ORDER BY f.name ASC LIMIT 5 OFFSET 0",
			values: []interface{}{"something"},
		},
	}

	for _, tst := range testResults {
		t.Run(tst.name, func(t *testing.T) {
			query := tst.query()
			buf := bytes.Buffer{}

			if err := query.Build(&buf); err != tst.err {
				t.Fatal(err)
			} else if buf.String() != tst.result {
				t.Fatalf("got: %s -- expected: %s", buf.String(), tst.result)
			} else if !reflect.DeepEqual(query.Params(), tst.values) {
				t.Fatalf("got: %v -- expected: %v", query.Params(), tst.values)
			}
		})
	}
}
