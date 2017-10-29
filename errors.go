package qb

import "errors"

var (
	ErrInvalidPointer = errors.New("qb: attempt to load into an invalid pointer")
)
