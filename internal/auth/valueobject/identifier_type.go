package valueobject

import (
	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	appErrors "github.com/teacinema-go/gateway-service/internal/errors"
	appValidator "github.com/teacinema-go/gateway-service/pkg/validator"
)

type IdentifierType string

const (
	IdentifierTypePhone IdentifierType = "phone"
	IdentifierTypeEmail IdentifierType = "email"
)

func (it IdentifierType) ToProto() authv1.IdentifierType {
	switch it {
	case IdentifierTypeEmail:
		return authv1.IdentifierType_EMAIL
	case IdentifierTypePhone:
		return authv1.IdentifierType_PHONE
	default:
		return authv1.IdentifierType_IDENTIFIER_TYPE_UNSPECIFIED
	}
}

func (it IdentifierType) Validate(identifier string) error {
	switch it {
	case IdentifierTypeEmail:
		if !appValidator.IsValidEmail(identifier) {
			return appErrors.ErrInvalidEmail
		}
	case IdentifierTypePhone:
		if !appValidator.IsValidE164Phone(identifier) {
			return appErrors.ErrInvalidE164Phone
		}
	default:
		return appErrors.ErrInvalidIdentifierType
	}

	return nil
}
