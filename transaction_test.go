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
	if err := tx.LoadValue(ctx, tx.Select().From("animals").Columns("COUNT(name)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 0 {
		t.Fatalf("Expected 0 record but got %d", totalCount)
	}

	// Insert one record
	if _, err := tx.Exec(ctx, tx.Insert().InTo("animals").Columns("name").Values("fuu")); err != nil {
		t.Fatal(err)
	}

	// Check one record is inserted
	if err := tx.LoadValue(ctx, tx.Select().From("animals").Columns("COUNT(name)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}

	// TODO Check if outsider still sees 0 records
	// if err := db.LoadValue(ctx, db.Select().From("animals").Columns("COUNT(name)"), &totalCount); err != nil {
	// 	t.Fatal(err)
	// } else if totalCount != 0 {
	// 	t.Fatalf("Expected 0 record but got %d", totalCount)
	// }

	// Insert a second record
	if _, err := tx.Exec(ctx, tx.Insert().InTo("animals").Columns("name").Values("bar")); err != nil {
		t.Fatal(err)
	}

	// Check two records are inserted
	if err := tx.LoadValue(ctx, tx.Select().From("animals").Columns("COUNT(name)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 2 {
		t.Fatalf("Expected 2 record but got %d", totalCount)
	}

	// Generate a conflict
	if _, err := tx.Exec(ctx, tx.Update().Table("animals").Set("name", "fuu").Where("name = ?", "bar")); err == nil {
		t.Fatalf("Expected unique constraint to kick in when UPDATE but it did not")
	} else if err.Error() != "constraint failed: UNIQUE constraint failed: animals.name (2067)" {
		t.Fatalf("got: %s -- expected: UNIQUE constraint failed: animals.name", err.Error())
	}

	// Rollback the transaction
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	// Make sure the transaction is closed
	if _, err := tx.Exec(ctx, tx.Delete().From("animals").Where("name = ?", "bar")); err == nil {
		t.Fatalf("Expected unique constraint to kick in when UPDATE but it did not")
	} else if err.Error() != "sql: transaction has already been committed or rolled back" {
		t.Fatalf("got: %s -- expected: sql: transaction has already been committed or rolled back", err.Error())
	}

	// Check if 0 records exist
	if err := db.LoadValue(ctx, db.Select().From("animals").Columns("COUNT(name)"), &totalCount); err != nil {
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
	if _, err := tx.Exec(ctx, tx.Insert().InTo("animals").Columns("name").Values("fuu")); err != nil {
		t.Fatal(err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Check if 1 record exists
	if err := db.LoadValue(ctx, db.Select().From("animals").Columns("COUNT(name)"), &totalCount); err != nil {
		t.Fatal(err)
	} else if totalCount != 1 {
		t.Fatalf("Expected 1 record but got %d", totalCount)
	}
}
