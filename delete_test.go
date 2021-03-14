package qb

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDeleteQuery(t *testing.T) {
	type test struct {
		name   string
		query  func() *DeleteQuery
		result string
		values []interface{}
		err    error
	}

	var testResults = []test{
		{
			name: "delete everything",
			query: func() *DeleteQuery {
				query := &DeleteQuery{table: "fuu"}
				return query
			},
			result: "DELETE FROM fuu",
		},
		{
			name: "delete with where clause",
			query: func() *DeleteQuery {
				query := &DeleteQuery{table: "fuu"}
				query.Where("column1 = ?", 1234)
				return query
			},
			result: "DELETE FROM fuu WHERE column1 = ?",
			values: []interface{}{1234},
		},
		{
			name: "delete with multiple where clauses",
			query: func() *DeleteQuery {
				query := &DeleteQuery{table: "fuu"}
				query.Where("column1 = ?", 1234)
				query.Where("column2 = ?", "test")
				query.Where("column3 IS NULL")
				return query
			},
			result: "DELETE FROM fuu WHERE column1 = ? AND column2 = ? AND column3 IS NULL",
			values: []interface{}{1234, "test"},
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
