package request

import (
	"github.com/go-playground/validator/v10"
	appErrors "github.com/teacinema-go/gateway-service/internal/errors"
	"github.com/teacinema-go/gateway-service/pkg/enum"
	appValidator "github.com/teacinema-go/gateway-service/pkg/validator"
)

type SendOtpRequest struct {
	Identifier     string              `json:"identifier" validate:"required"`
	IdentifierType enum.IdentifierType `json:"identifier_type" validate:"required,oneof=phone email"`
}

func (r *SendOtpRequest) Validate(v *validator.Validate) error {
	// TODO format validation messages
	if err := v.Struct(r); err != nil {
		return err
	}

	return validateIdentifier(r.Identifier, r.IdentifierType)
}

type VerifyOtpRequest struct {
	Identifier     string              `json:"identifier" validate:"required"`
	IdentifierType enum.IdentifierType `json:"identifier_type" validate:"required,oneof=phone email"`
	Otp            string              `json:"otp" validate:"required"`
}

func (r *VerifyOtpRequest) Validate(v *validator.Validate) error {
	// TODO format validation messages
	if err := v.Struct(r); err != nil {
		return err
	}

	return validateIdentifier(r.Identifier, r.IdentifierType)
}

func validateIdentifier(identifier string, identifierType enum.IdentifierType) error {
	switch identifierType {
	case enum.IdentifierTypeEmail:
		if !appValidator.IsValidEmail(identifier) {
			return appErrors.ErrInvalidEmail
		}
	case enum.IdentifierTypePhone:
		if !appValidator.IsValidE164Phone(identifier) {
			return appErrors.ErrInvalidE164Phone
		}
	default:
		return appErrors.ErrInvalidIdentifierType
	}

	return nil
}
