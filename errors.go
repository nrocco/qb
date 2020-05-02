package qb

import "errors"

var (
	// ErrInvalidPointer indicates that you passed an invalid pointer into a function
	ErrInvalidPointer = errors.New("qb: attempt to load into an invalid pointer")
)
