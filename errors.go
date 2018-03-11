package srv

import "errors"

var (
	// ErrConflict is the conflict error.
	ErrConflict = errors.New("conflict error")
	// ErrNotFound is the not found error.
	ErrNotFound = errors.New("not found")
)
