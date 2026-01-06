package dto

import (
	"github.com/go-playground/validator/v10"
	appErrors "github.com/teacinema-go/gateway-service/internal/errors"
	appValidator "github.com/teacinema-go/gateway-service/pkg/validator"
)

type IdentifierType string

const (
	IdentifierPhone IdentifierType = "phone"
	IdentifierEmail IdentifierType = "email"
)

type SendOtpRequest struct {
	Identifier string         `json:"identifier" validate:"required"`
	Type       IdentifierType `json:"type" validate:"required,oneof=phone email"`
}

func (r *SendOtpRequest) Validate(v *validator.Validate) error {
	// TODO format validation messages
	if err := v.Struct(r); err != nil {
		return err
	}

	switch r.Type {
	case IdentifierEmail:
		if !appValidator.IsValidEmail(r.Identifier) {
			return appErrors.InvalidEmailError("identifier")
		}
	case IdentifierPhone:
		if !appValidator.IsValidE164Phone(r.Identifier) {
			return appErrors.InvalidE164PhoneError("identifier")
		}
	default:
		return appErrors.InvalidFieldError("identifier type")
	}

	return nil
}
