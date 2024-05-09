package model

import "errors"

var (
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
)
