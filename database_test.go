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
	query := db.For(ctx).Select().From("notes").Columns("COUNT(id)")

	totalCount := 0
	if err := query.LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	query.Columns("*")

	notes := []*note{}
	if _, err := query.Load(ctx, &notes); err != nil {
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

	query := db.For(ctx).Insert().InTo("notes").Columns("name", "content").Record(&n).WriteID(&n.ID)

	if _, err := query.Exec(ctx); err != nil {
		t.Fatal(err)
	}

	if n.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", n.ID)
	}
}

func TestUpdateIntoDatabase(t *testing.T) {
	db := createTestDB(t, notesSchema, `INSERT INTO notes (id, name, content) VALUES (1, "Fuu", "This is bar");`)
	defer db.Close()

	ctx := context.TODO()
	n := note{}

	preUpdateQuery := db.For(ctx).Select().From("notes").Where("id = ?", 1)
	if _, err := preUpdateQuery.Load(ctx, &n); err != nil {
		t.Fatal(err)
	} else if n.Name != "Fuu" {
		t.Fatalf("Expected Fuu but got %s", n.Name)
	} else if n.ID != 1 {
		t.Fatalf("Expected note.ID to be 1 but got %d", n.ID)
	}

	updateQuery := db.For(ctx).Update().Table("notes").Set("name", "Bar").Where("id = ?", n.ID)
	if _, err := updateQuery.Exec(ctx); err != nil {
		t.Fatal(err)
	}

	postUpdateQuery := db.For(ctx).Select().From("notes")
	if _, err := postUpdateQuery.Load(ctx, &n); err != nil {
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

	preDeleteQuery := db.For(ctx).Select().From("notes").Columns("COUNT(id)")
	if err := preDeleteQuery.LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}

	deleteQuery := db.For(ctx).Delete().From("notes").Where("id = ?", 1)
	if _, err := deleteQuery.Exec(ctx); err != nil {
		t.Fatal(err)
	}

	postDeleteQuery := db.For(ctx).Select().From("notes").Columns("COUNT(id)")
	if err := postDeleteQuery.LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 0 {
		t.Fatalf("Expected 0 record but got %d", totalCount)
	}
}
