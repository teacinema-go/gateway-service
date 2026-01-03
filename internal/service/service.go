package service

import (
	"context"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/dto"
)

type AuthService interface {
	SendOtp(ctx context.Context, req dto.SendOtpRequest) (*authv1.SendOtpResponse, error)
	Close() error
}
