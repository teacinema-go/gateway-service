package errors

import "errors"

var (
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPhone          = errors.New("invalid phone format (e164)")
	ErrInvalidIdentifierType = errors.New("invalid identifier type")
)
