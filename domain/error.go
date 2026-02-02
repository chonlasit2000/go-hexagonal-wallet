package domain

import "errors"

var (
	// General errors
	ErrBadRequest             = errors.New("bad request")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrNotFound               = errors.New("not found")
	ErrConflict               = errors.New("conflict")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInternalServerError    = errors.New("internal server error")

	// Wallet specific errors
	ErrInsufficientBalance = errors.New("insufficient balance")    // เงินไม่พอ
	ErrSameAccount         = errors.New("cannot transfer to self") // โอนให้ตัวเองไม่ได้
)
