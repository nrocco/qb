package qb

import (
	"context"
	"testing"
)

const animalsSchema = `CREATE TABLE animals (name TEXT NOT NULL, UNIQUE(name));`

func TestTransactionsRollback(t *testing.T) {
	db := createTestDB(t, animalsSchema, "")
	defer db.Close()

	totalCount := 0
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure we start with an empty table
	if err := tx.Select("animals").Columns("COUNT(name)").LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 0 {
		t.Fatalf("Expected 0 record but got %d", totalCount)
	}

	// Insert one record
	if _, err := tx.Insert("animals").Columns("name").Values("fuu").Exec(ctx); err != nil {
		t.Fatal(err)
	}

	// Check one record is inserted
	if err := tx.Select("animals").Columns("COUNT(name)").LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}

	// TODO Check if outsider still sees 0 records
	// if err := db.Select("animals").Columns("COUNT(name)").LoadValue(ctx, &totalCount); err != nil {
	// 	t.Fatal(err)
	// } else if totalCount != 0 {
	// 	t.Fatalf("Expected 0 record but got %d", totalCount)
	// }

	// Insert a second record
	if _, err := tx.Insert("animals").Columns("name").Values("bar").Exec(ctx); err != nil {
		t.Fatal(err)
	}

	// Check two records are inserted
	if err := tx.Select("animals").Columns("COUNT(name)").LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	// Generate a conflict
	if _, err := tx.Update("animals").Set("name", "fuu").Where("name = ?", "bar").Exec(ctx); err == nil {
		t.Fatalf("Expected unique constraint to kick in when UPDATE but it did not")
	} else if err.Error() != "UNIQUE constraint failed: animals.name" {
		t.Fatalf("got: %s -- expected: UNIQUE constraint failed: animals.name", err.Error())
	}

	// Rollback the transaction
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	// Make sure the transaction is closed
	if _, err := tx.Delete("animals").Where("name = ?", "bar").Exec(ctx); err == nil {
		t.Fatalf("Expected unique constraint to kick in when UPDATE but it did not")
	} else if err.Error() != "sql: transaction has already been committed or rolled back" {
		t.Fatalf("got: %s -- expected: sql: transaction has already been committed or rolled back", err.Error())
	}

	// Check if 0 records exist
	if err := db.Select("animals").Columns("COUNT(name)").LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 0 {
		t.Fatalf("Expected 0 record but got %d", totalCount)
	}
}

func TestTransactionsCommited(t *testing.T) {
	db := createTestDB(t, animalsSchema, "")
	defer db.Close()

	totalCount := 0
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Insert one record
	if _, err := tx.Insert("animals").Columns("name").Values("fuu").Exec(ctx); err != nil {
		t.Fatal(err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Check if 0 records exist
	if err := db.Select("animals").Columns("COUNT(name)").LoadValue(ctx, &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}
}
