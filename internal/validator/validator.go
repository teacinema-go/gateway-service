package validator

import (
	"regexp"

	appErrors "github.com/teacinema-go/gateway-service/internal/errors"

	"github.com/nyaruka/phonenumbers"
	"github.com/teacinema-go/gateway-service/internal/dto"
)

func ValidateIdentifierFormat(req dto.SendOtpRequest) error {
	switch req.Type {
	case dto.IdentifierEmail:
		if !isValidEmail(req.Identifier) {
			return appErrors.InvalidEmailError("identifier")
		}
	case dto.IdentifierPhone:
		if !isValidE164Phone(req.Identifier) {
			return appErrors.InvalidE164PhoneError("identifier")
		}
	default:
		return appErrors.InvalidFieldError("identifier type")
	}

	return nil
}

func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidE164Phone(phone string) bool {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false
	}

	if !phonenumbers.IsValidNumber(num) {
		return false
	}

	formatted := phonenumbers.Format(num, phonenumbers.E164)
	return formatted == phone
}
