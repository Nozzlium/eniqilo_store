package constant

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("account already exists")
	ErrBadInput = errors.New("invalid input")
)
