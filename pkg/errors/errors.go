package errors

import "errors"

var (
	ErrConflict = errors.New("conflict error")
	ErrNotFound = errors.New("not found")
)
