package constant

import "errors"

var (
	ErrBadInput      = errors.New("invalid input")
	ErrProductExists = errors.New("product already exists")

	ErrConflict = errors.New(
		"account already exists",
	)

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
