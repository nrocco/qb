package qb

import (
	"bytes"
	"reflect"
	"testing"
)

func TestInsertQuery(t *testing.T) {
	type test struct {
		name   string
		query  func() *InsertQuery
		result string
		values []interface{}
		err    error
	}

	var testResults = []test{
		test{
			name: "insert nothing", // TODO consider throwing an error here
			query: func() *InsertQuery {
				query := &InsertQuery{table: "fuu"}
				return query
			},
			result: "INSERT INTO fuu VALUES ()",
		},
		test{
			name: "insert and returning",
			query: func() *InsertQuery {
				query := &InsertQuery{table: "fuu"}
				query.Values(123)
				query.Returning("column1", "column3")
				return query
			},
			result: "INSERT INTO fuu VALUES (?) RETURNING column1, column3",
			values: []interface{}{123},
		},
		test{
			name: "insert or ignore",
			query: func() *InsertQuery {
				query := &InsertQuery{table: "fuu"}
				query.OrIgnore()
				query.Values(123, "fuubar")
				return query
			},
			result: "INSERT OR IGNORE INTO fuu VALUES (?, ?)",
			values: []interface{}{123, "fuubar"},
		},
		test{
			name: "insert one column",
			query: func() *InsertQuery {
				query := &InsertQuery{table: "fuu"}
				query.Columns("column1")
				query.Values("value1")
				return query
			},
			result: "INSERT INTO fuu (column1) VALUES (?)",
			values: []interface{}{"value1"},
		},
		test{
			name: "insert with on conflict",
			query: func() *InsertQuery {
				query := &InsertQuery{table: "fuu"}
				query.Columns("column1")
				query.Values("value1")
				query.OnConflict("column1", "column1=excluded.column1, column2=excluded.column2")
				return query
			},
			result: "INSERT INTO fuu (column1) VALUES (?) ON CONFLICT (column1) DO UPDATE SET column1=excluded.column1, column2=excluded.column2",
			values: []interface{}{"value1"},
		},
		test{
			name: "insert record",
			query: func() *InsertQuery {
				record := struct {
					ID      int64
					Name    string
					Content string
					FuuBar  string
					NonExistent int
				}{
					ID:      12345,
					Name:    "fuubar",
					Content: "something else",
					FuuBar:  "testtest",
				}
				query := &InsertQuery{table: "fuu"}
				query.Columns("id", "name", "fuu_bar", "non_existent")
				query.Record(&record)
				return query
			},
			result: "INSERT INTO fuu (id, name, fuu_bar, non_existent) VALUES (?, ?, ?, ?)",
			values: []interface{}{int64(12345), "fuubar", "testtest", 0},
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
