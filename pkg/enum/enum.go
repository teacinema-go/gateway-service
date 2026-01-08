package enum

import (
	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
)

type IdentifierType string

const (
	IdentifierTypePhone IdentifierType = "phone"
	IdentifierTypeEmail IdentifierType = "email"
)

func (e IdentifierType) ToProto() authv1.IdentifierType {
	identifierType := authv1.IdentifierType_IDENTIFIER_TYPE_UNSPECIFIED
	if e == IdentifierTypeEmail {
		identifierType = authv1.IdentifierType_EMAIL
	}

	if e == IdentifierTypePhone {
		identifierType = authv1.IdentifierType_PHONE
	}

	return identifierType
}
