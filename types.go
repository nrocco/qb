package qb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

var nullString = []byte("null")

type NullString struct {
	sql.NullString
}

// MarshalJSON correctly serializes a NullString to JSON
func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a NullString from JSON
func (n *NullString) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

var nullInt64 = []byte("null")

type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON correctly serializes a NullInt64 to JSON
func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return nullInt64, nil
}

// UnmarshalJSON correctly deserializes a NullInt64 from JSON
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

var nullTime = []byte("null")

type NullTime struct {
	time.Time
	Valid bool
}

// Scan a raw value and wrap it in NullTime
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value returns the underlying value Time
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON correctly serializes a NullTime to JSON
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return nullTime, nil
}

// UnmarshalJSON correctly deserializes a NullTime from JSON
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		return nt.Scan(s)
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	return nt.Scan(t)
}
