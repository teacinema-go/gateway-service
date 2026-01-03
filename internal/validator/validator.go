package validator

import (
	appErrors "github.com/teacinema-go/gateway-service/internal/errors"
	"github.com/teacinema-go/gateway-service/pkg/validator"

	"github.com/teacinema-go/gateway-service/internal/dto"
)

func ValidateIdentifierFormat(req dto.SendOtpRequest) error {
	switch req.Type {
	case dto.IdentifierEmail:
		if !validator.IsValidEmail(req.Identifier) {
			return appErrors.InvalidEmailError("identifier")
		}
	case dto.IdentifierPhone:
		if !validator.IsValidE164Phone(req.Identifier) {
			return appErrors.InvalidE164PhoneError("identifier")
		}
	default:
		return appErrors.InvalidFieldError("identifier type")
	}

	return nil
}
