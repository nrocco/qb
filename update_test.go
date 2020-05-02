package qb

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUpdateQuery(t *testing.T) {
	type test struct {
		name   string
		query  func() *UpdateQuery
		result string
		values []interface{}
		err    error
	}

	var testResults = []test{
		test{
			name: "update nothing", // TODO consider throwing an error here
			query: func() *UpdateQuery {
				query := &UpdateQuery{table: "fuu"}
				return query
			},
			result: "UPDATE fuu SET ",
		},
		test{
			name: "update and returning",
			query: func() *UpdateQuery {
				query := &UpdateQuery{table: "fuu"}
				query.Set("closed", true)
				query.Where("id = ?", 123)
				query.Returning("column1", "column3")
				return query
			},
			result: "UPDATE fuu SET closed = ? WHERE id = ? RETURNING column1, column3",
			values: []interface{}{true, 123},
		},
		test{
			name: "update one column",
			query: func() *UpdateQuery {
				query := &UpdateQuery{table: "fuu"}
				query.Set("closed", true)
				query.Where("id = ?", 123)
				return query
			},
			result: "UPDATE fuu SET closed = ? WHERE id = ?",
			values: []interface{}{true, 123},
		},
		test{
			name: "update multiple columns multiple where clauses",
			query: func() *UpdateQuery {
				query := &UpdateQuery{table: "fuu"}
				query.Set("closed", true)
				query.Set("year", 2020)
				query.Where("id = ?", 123)
				query.Where("name = ?", "test")
				return query
			},
			result: "UPDATE fuu SET closed = ?, year = ? WHERE id = ? AND name = ?",
			values: []interface{}{true, 2020, 123, "test"},
		},
		test{
			name: "update without where clause",
			query: func() *UpdateQuery {
				query := &UpdateQuery{table: "fuu"}
				query.Set("closed", true)
				return query
			},
			result: "UPDATE fuu SET closed = ?",
			values: []interface{}{true},
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
