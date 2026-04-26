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

	ctx := context.TODO()

	totalCount := 0
	if err := db.LoadValue(ctx, db.Select().From("names").Columns("COUNT(name)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	names := []*name{}
	if _, err := db.Load(ctx, db.Select().From("names").Columns("*"), &names); err != nil {
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

	ctx := context.TODO()

	n := name{
		Name: "Test Name",
		Tags: Tags{"tag1", "tag2"},
	}

	if _, err := db.Exec(ctx, db.Insert().InTo("names").Columns("name", "tags").Record(&n)); err != nil {
		t.Fatal(err)
	}

	fuu := name{}
	if _, err := db.Load(ctx, db.Select().From("names").Columns("*").Where("name = ?", "Test Name"), &fuu); err != nil {
		t.Fatal(err)
	}

	if len(fuu.Tags) != 2 {
		t.Fatalf("Expected 2 tags but got `%d`: %v", len(fuu.Tags), fuu.Tags)
	}

	if (fuu.Tags)[0] != "tag1" {
		t.Fatalf("Expected tag `tag1` but got `%s`", (fuu.Tags)[0])
	}
}
