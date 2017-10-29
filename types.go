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
