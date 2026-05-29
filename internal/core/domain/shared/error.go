package shared

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrAlreadyExist = errors.New("already exists")
	ErrInvalidInput = errors.New("invalid input")
)
