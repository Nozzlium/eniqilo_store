package constant

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("account already exists")
	ErrBadInput      = errors.New("invalid input")
	ErrProductExists = errors.New("product already exists")

	ErrInvalidBody = errors.New(
		"invalid body",
	)

	ErrCannotEdit = errors.New(
		"cannot edit",
	)

	ErrSavingData = errors.New(
		"failed to save data",
	)
)
