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
	ID      int64
	Name    string
	Content string
}

func createTestDB(t *testing.T, schema string, fixtures string) *DB {
	db, err := Open(context.TODO(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if schema != "" {
		if _, err = db.DB.Exec(schema); err != nil {
			t.Fatalf("Could not create test schema: %s", err)
		}
	}

	if fixtures != "" {
		if _, err = db.DB.Exec(fixtures); err != nil {
			t.Fatalf("Could not load fixtures: %s", err)
		}
	}

	return db
}

func TestOpenDatabase(t *testing.T) {
	_, err := Open(context.TODO(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSelectFromDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES
	(1, "Fuu", "This is bar"),
	(2, "Test", "This is fuu")`)

	defer db.Close()

	ctx := context.TODO()
	q := db.Select().From("notes").Columns("COUNT(id)")

	totalCount := 0
	if err := db.LoadValue(ctx, q, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	q.Columns("*")

	notes := []*note{}
	if _, err := db.Load(ctx, q, &notes); err != nil {
		t.Fatal(err)
	} else if len(notes) != 2 {
		t.Fatalf("Expected 2 rows but got %d", len(notes))
	}
}

func TestInsertIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, "")
	defer db.Close()

	ctx := context.TODO()

	n := note{
		Name:    "Test Name",
		Content: "Test Content",
	}

	q := db.Insert().InTo("notes").Columns("name", "content").Record(&n)

	result, err := db.Exec(ctx, q)
	if err != nil {
		t.Fatal(err)
	}

	n.ID, _ = result.LastInsertId()

	if n.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", n.ID)
	}
}

func TestUpdateIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES (1, "Fuu", "This is bar");`)
	defer db.Close()

	ctx := context.TODO()
	n := note{}

	if _, err := db.Load(ctx, db.Select().From("notes").Where("id = ?", 1), &n); err != nil {
		t.Fatal(err)
	} else if n.Name != "Fuu" {
		t.Fatalf("Expected Fuu but got %s", n.Name)
	} else if n.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", n.ID)
	}

	if _, err := db.Exec(ctx, db.Update().Table("notes").Set("name", "Bar").Where("id = ?", n.ID)); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Load(ctx, db.Select().From("notes"), &n); err != nil {
		t.Fatal(err)
	} else if n.Name != "Bar" {
		t.Fatalf("Expected Bar but got %s", n.Name)
	} else if n.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", n.ID)
	}
}

func TestDeleteIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES (1, "Fuu", "This is bar");`)
	defer db.Close()

	ctx := context.TODO()
	totalCount := 0

	if err := db.LoadValue(ctx, db.Select().From("notes").Columns("COUNT(id)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}

	if _, err := db.Exec(ctx, db.Delete().From("notes").Where("id = ?", 1)); err != nil {
		t.Fatal(err)
	}

	if err := db.LoadValue(ctx, db.Select().From("notes").Columns("COUNT(id)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 0 {
		t.Fatalf("Expected 0 record but got %d", totalCount)
	}
}
