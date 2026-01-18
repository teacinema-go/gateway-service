package clients

import (
	"context"
	"fmt"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/auth/dto/request"
	"github.com/teacinema-go/gateway-service/internal/config"
)

type AuthServiceClient interface {
	SendOtp(ctx context.Context, req request.SendOtpRequest) (*authv1.SendOtpResponse, error)
	VerifyOtp(ctx context.Context, req request.VerifyOtpRequest) (*authv1.VerifyOtpResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*authv1.RefreshResponse, error)
	Close() error
}

type Manager struct {
	Auth    AuthServiceClient
	clients []Closer
}

type Closer interface {
	Close() error
}

func NewClientManager(serviceConfig config.ServiceConfig) (*Manager, error) {
	authService, err := NewAuthServiceClient(serviceConfig.AuthServiceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	return &Manager{
		Auth:    authService,
		clients: []Closer{authService},
	}, nil
}

func (m *Manager) Close() error {
	var firstErr error
	for _, client := range m.clients {
		if err := client.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
