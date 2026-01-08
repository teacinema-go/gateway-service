package errors

import (
	"errors"
)

var (
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidE164Phone      = errors.New("invalid e.164 phone number")
	ErrInvalidIdentifierType = errors.New("invalid identifier type")
)
