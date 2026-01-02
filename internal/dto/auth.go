package dto

type IdentifierType string

const (
	IdentifierPhone IdentifierType = "phone"
	IdentifierEmail IdentifierType = "email"
)

type SendOtpRequest struct {
	Identifier string         `json:"identifier" validate:"required"`
	Type       IdentifierType `json:"type" validate:"required,oneof=phone email"`
}
