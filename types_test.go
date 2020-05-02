package qb

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func TestNullStringMarshalInvalid(t *testing.T) {
	value := NullString{}

	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(result, []byte("null")) {
		t.Fatalf("Expected null but got %s", result)
	}
}

func TestNullStringUnmarshalInvalid(t *testing.T) {
	value := NullString{}
	result := []byte("null")

	err := json.Unmarshal(result, &value)
	if err != nil {
		t.Fatal(err)
	} else if value.String != "" {
		t.Fatalf("Expected empty string but got %s", value.String)
	} else if value.Valid {
		t.Fatal("Expected NullString to not be valid")
	}
}

func TestNullStringMarshalValid(t *testing.T) {
	value := NullString{sql.NullString{"test", true}}
	expected := []byte("\"test\"")

	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestNullStringUnmarshalValid(t *testing.T) {
	value := NullString{}
	result := []byte("\"test\"")
	expected := "test"

	err := json.Unmarshal(result, &value)
	if err != nil {
		t.Fatal(err)
	} else if value.String != expected {
		t.Fatalf("Expected %s but got %s", expected, value.String)
	} else if !value.Valid {
		t.Fatal("Expected NullString to be valid")
	}
}

func TestNullInt64MarshalInvalid(t *testing.T) {
	value := NullInt64{}

	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(result, []byte("null")) {
		t.Fatalf("Expected null but got %s", result)
	}
}

func TestNullInt64UnmarshalInvalid(t *testing.T) {
	value := NullInt64{}
	result := []byte("null")

	err := json.Unmarshal(result, &value)
	if err != nil {
		t.Fatal(err)
	} else if value.Int64 != 0 {
		t.Fatalf("Expected 0 but got %d", value.Int64)
	} else if value.Valid {
		t.Fatal("Expected NullInt64 to not be valid")
	}
}

func TestNullInt64MarshalValid(t *testing.T) {
	value := NullInt64{sql.NullInt64{123, true}}
	expected := []byte("123")

	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(result, expected) {
		t.Fatalf("Expected %d but got %d", expected, result)
	}
}

func TestNullInt64UnmarshalValid(t *testing.T) {
	value := NullInt64{}
	result := []byte("123")
	expected := int64(123)

	err := json.Unmarshal(result, &value)
	if err != nil {
		t.Fatal(err)
	} else if value.Int64 != expected {
		t.Fatalf("Expected %d but got %d", expected, value.Int64)
	} else if !value.Valid {
		t.Fatal("Expected NullInt64 to be valid")
	}
}

func TestNullTimeMarshalInvalid(t *testing.T) {
	value := NullTime{}

	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(result, []byte("null")) {
		t.Fatalf("Expected null but got %s", result)
	}
}

func TestNullTimeUnmarshalInvalid(t *testing.T) {
	value := NullTime{}
	result := []byte("null")

	err := json.Unmarshal(result, &value)
	if err != nil {
		t.Fatal(err)
	} else if !value.Time.IsZero() {
		t.Fatalf("Expected empty string but got %v", value.Time)
	} else if value.Valid {
		t.Fatal("Expected NullTime to not be valid")
	}
}

func TestNullTimeValid(t *testing.T) {
	value := NullTime{sql.NullTime{time.Date(2015, 9, 18, 0, 0, 0, 0, time.UTC), true}}
	expected := []byte("\"2015-09-18T00:00:00Z\"")

	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, expected) {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestNullTimeUnmarshalValid(t *testing.T) {
	value := NullTime{}
	result := []byte("\"2015-09-18T00:00:00Z\"")
	expected := time.Date(2015, 9, 18, 0, 0, 0, 0, time.UTC)

	err := json.Unmarshal(result, &value)
	if err != nil {
		t.Fatal(err)
	} else if !expected.Equal(value.Time) {
		t.Fatalf("Expected %v but got %v", expected, value.Time)
	} else if !value.Valid {
		t.Fatal("Expected NullTime to be valid")
	}
}
