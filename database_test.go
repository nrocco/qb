package qb

import (
	"context"
	"testing"
)

const notesSchema = `CREATE TABLE notes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(64) NOT NULL UNIQUE,
	content VARCHAR(255) NULL
);`

type note struct {
	ID int64
	Name string
	Content string
}

func createTestDB(t *testing.T, schema string, fixtures string) *DB {
	db, err := Open(context.TODO(), ":memory:", nil)
	if err != nil {
		t.Fatal(err)
	}

	if schema != "" {
		if _, err = db.Exec(schema); err != nil {
			t.Fatalf("Could not create test schema: %s", err)
		}
	}

	if fixtures != "" {
		if _, err = db.Exec(fixtures); err != nil {
			t.Fatalf("Could not load fixtures: %s", err)
		}
	}

	return db
}

func TestOpenDatabase(t *testing.T) {
	db, err := Open(context.TODO(), ":memory:", nil)
	if err != nil {
		t.Fatal(err)
	}

	if db.logger != nil {
		t.Fatal("Expected logger to be nil")
	}
}

func TestSelectFromDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES
	(1, "Fuu", "This is bar"),
	(2, "Test", "This is fuu")`)

	defer db.Close()

	query := db.Select("notes")
	query.Columns("COUNT(id)")

	totalCount := 0
	if err := query.LoadValue(&totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	query.Columns("*")

	notes := []*note{}
	if _, err := query.Load(&notes); err != nil {
		t.Fatal(err)
	} else if len(notes) != 2 {
		t.Fatalf("Expected 2 rows but got %d", len(notes))
	}
}

func TestInsertIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, "")
	defer db.Close()

	note := note{
		Name: "Test Name",
		Content: "Test Content",
	}

	query := db.Insert("notes")
	query.Columns("name", "content")
	query.Record(&note)

	if _, err := query.Exec(); err != nil {
		t.Fatal(err)
	}

	if note.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", note.ID)
	}
}

func TestUpdateIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES (1, "Fuu", "This is bar");`)
	defer db.Close()

	note := note{}
	preUpdateQuery := db.Select("notes").Where("id = ?", 1)
	if _, err := preUpdateQuery.Load(&note); err != nil {
		t.Fatal(err)
	} else if note.Name != "Fuu" {
		t.Fatalf("Expected Fuu but got %s", note.Name)
	} else if note.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", note.ID)
	}

	updateQuery := db.Update("notes").Set("name", "Bar").Where("id = ?", note.ID)
	if _, err := updateQuery.Exec(); err != nil {
		t.Fatal(err)
	}

	postUpdateQuery := db.Select("notes")
	if _, err := postUpdateQuery.Load(&note); err != nil {
		t.Fatal(err)
	} else if note.Name != "Bar" {
		t.Fatalf("Expected Bar but got %s", note.Name)
	} else if note.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", note.ID)
	}
}

func TestDeleteIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES (1, "Fuu", "This is bar");`)
	defer db.Close()

	totalCount := 0
	preDeleteQuery := db.Select("notes").Columns("COUNT(id)")
	if err := preDeleteQuery.LoadValue(&totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}

	deleteQuery := db.Delete("notes").Where("id = ?", 1)
	if _, err := deleteQuery.Exec(); err != nil {
		t.Fatal(err)
	}

	postDeleteQuery := db.Select("notes").Columns("COUNT(id)")
	if err := postDeleteQuery.LoadValue(&totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 0 {
		t.Fatalf("Expected 0 record but got %d", totalCount)
	}
}
