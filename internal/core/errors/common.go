package core_errors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrConflict            = errors.New("conflict")
	ErrInternalServerError = errors.New("internal server error")
)
