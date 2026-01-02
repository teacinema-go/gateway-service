package errors

import (
	"fmt"
)

func InvalidEmailError(field string) error {
	return fmt.Errorf("%s must be a valid email", field)
}

func InvalidE164PhoneError(field string) error {
	return fmt.Errorf("%s must be a valid e.164 phone number", field)
}

func InvalidFieldError(field string) error {
	return fmt.Errorf("%s is invalid", field)
}
