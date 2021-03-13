package qb

import (
	"context"
	"testing"
)

const fuuSchema = `PRAGMA foreign_keys = ON;

CREATE TABLE artist(
  id INTEGER PRIMARY KEY,
  name TEXT
);

INSERT INTO artist VALUES (1, 'Dean Martin');
INSERT INTO artist VALUES (2, 'Frank Sinatra');

CREATE TABLE track(
  id INTEGER,
  name TEXT,
  artist INTEGER,
  FOREIGN KEY(artist) REFERENCES artist(id)
);

INSERT INTO track VALUES (11, 'That is Amore', 1);
INSERT INTO track VALUES (12, 'Christmas Blues', 1);
INSERT INTO track VALUES (13, 'My Way', 2);
`

func TestForeignKeys(t *testing.T) {
	db := createTestDB(t, fuuSchema, "")
	defer db.Close()

	totalCount := 0
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Insert one record
	if _, err := tx.Insert(ctx).InTo("track").Columns("id", "name", "artist").Values(14, "Fuubar", 3).Exec(); err == nil {
		t.Fatal("Expected a foreign key constraint error")
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Check if 0 records exist
	if err := db.Select(ctx).From("track").Columns("COUNT(id)").LoadValue(&totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 3 {
		t.Fatalf("Expected 3 record but got %d", totalCount)
	}
}
