package qb

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONValue is a convenient function to implement the Valuer interface for json columns
func JSONValue(value interface{}) (driver.Value, error) {
	v, err := json.Marshal(value)
	return string(v), err
}

// JSONScan is a convenient function to implement the Scanner interface for json columns
func JSONScan(destination interface{}, value interface{}) error {
	var err error

	switch v := value.(type) {
	case []uint8:
		err = json.Unmarshal(v, destination)
	case string:
		err = json.Unmarshal([]byte(v), destination)
	}

	return err
}
