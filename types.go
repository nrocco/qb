package qb

import (
	"database/sql"
	"encoding/json"
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
