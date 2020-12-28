package qb

import (
	"context"
	"database/sql/driver"
	"testing"
)

const namesSchema = `CREATE TABLE names (name TEXT NOT NULL, tags JSON NOT NULL);`

type Tags []string

func (t Tags) Value() (driver.Value, error) {
	return JSONValue(t)
}

func (t *Tags) Scan(value interface{}) error {
	return JSONScan(t, value)
}

type name struct {
	Name string
	Tags Tags
}

func TestSelectJSONFromDatabase(t *testing.T) {
	db := createTestDB(t, namesSchema, `INSERT INTO names (name, tags) VALUES ('Nico', '["fuu", "bar"]'), ('Tana', '["bar", "baz"]')`)

	defer db.Close()

	query := db.Select(context.TODO()).From("names")

	query.Columns("COUNT(name)")
	totalCount := 0
	if err := query.LoadValue(&totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	query.Columns("*")
	names := []*name{}
	if _, err := query.Load(&names); err != nil {
		t.Fatal(err)
	} else if len(names) != 2 {
		t.Fatalf("Expected 2 rows but got %d", len(names))
	}

	if len(names[0].Tags) != 2 {
		t.Fatalf("Expected 2 tags but got %d", len(names[0].Tags))
	} else if names[0].Tags[0] != "fuu" {
		t.Fatalf("Expected `fuu` but got `%s`", names[0].Tags[0])
	} else if names[0].Tags[1] != "bar" {
		t.Fatalf("Expected `fuu` but got `%s`", names[0].Tags[1])
	}

	if len(names[1].Tags) != 2 {
		t.Fatalf("Expected 2 tags but got %d", len(names[0].Tags))
	} else if names[1].Tags[0] != "bar" {
		t.Fatalf("Expected `fuu` but got `%s`", names[1].Tags[0])
	} else if names[1].Tags[1] != "baz" {
		t.Fatalf("Expected `fuu` but got `%s`", names[1].Tags[1])
	}
}

func TestInsertJSONIntoDatabase(t *testing.T) {
	db := createTestDB(t, namesSchema, "")
	defer db.Close()

	n := name{
		Name: "Test Name",
		Tags: Tags{
			"tag1",
			"tag2",
		},
	}

	query := db.Insert(context.TODO()).InTo("names")
	query.Columns("name", "tags")
	query.Record(&n)

	if _, err := query.Exec(); err != nil {
		t.Fatal(err)
	}

	iQuery := db.Select(context.TODO()).From("names")
	iQuery.Columns("*")
	iQuery.Where("name = ?", "Test Name")

	fuu := name{}

	if _, err := iQuery.Load(&fuu); err != nil {
		t.Fatal(err)
	}

	if len(fuu.Tags) != 2 {
		t.Fatalf("Expected 2 tags but got `%d`: %v", len(fuu.Tags), fuu.Tags)
	}

	if (fuu.Tags)[0] != "tag1" {
		t.Fatalf("Expected tag `tag1` but got `%s`", (fuu.Tags)[0])
	}
}
