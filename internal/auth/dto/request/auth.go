package request

import (
	"github.com/teacinema-go/gateway-service/internal/auth/valueobject"
)

type SendOtpRequest struct {
	Identifier     string                     `json:"identifier" validate:"required"`
	IdentifierType valueobject.IdentifierType `json:"identifier_type" validate:"required,oneof=phone email"`
}

type VerifyOtpRequest struct {
	Identifier     string                     `json:"identifier" validate:"required"`
	IdentifierType valueobject.IdentifierType `json:"identifier_type" validate:"required,oneof=phone email"`
	Otp            string                     `json:"otp" validate:"required,numeric,len=6"`
}

type UserID string
