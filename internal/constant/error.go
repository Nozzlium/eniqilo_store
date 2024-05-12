package constant

import "errors"

var (
	ErrConflict      = errors.New("account already exists")
	ErrBadInput      = errors.New("invalid input")
	ErrProductExists = errors.New("product already exists")

	ErrInvalidBody = errors.New(
		"invalid body",
	)

	ErrCannotEdit = errors.New(
		"cannot edit",
	)

	ErrNotFound = errors.New(
		"not found",
	)

	ErrSavingData = errors.New(
		"failed to save data",
	)

	ErrInsufficientFund = errors.New(
		"insufficient fund",
	)

	ErrInsufficientStock = errors.New(
		"insufficient stock or unavailable",
	)

	ErrInvalidChange = errors.New(
		"invalid change",
	)
)
