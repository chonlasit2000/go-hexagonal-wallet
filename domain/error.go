package domain

import "errors"

var (
	ErrBadRequest             = errors.New("bad request")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrNotFound               = errors.New("not found")
	ErrConflict               = errors.New("conflict")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInternalServerError    = errors.New("internal server error")
)
