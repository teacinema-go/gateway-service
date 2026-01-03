package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidE164Phone = errors.New("invalid e.164 phone number")
	ErrInvalidField     = errors.New("invalid field")
)

func InvalidEmailError(field string) error {
	return fmt.Errorf("%s: %w", field, ErrInvalidEmail)
}

func InvalidE164PhoneError(field string) error {
	return fmt.Errorf("%s: %w", field, ErrInvalidE164Phone)
}

func InvalidFieldError(field string) error {
	return fmt.Errorf("%s: %w", field, ErrInvalidField)
}
