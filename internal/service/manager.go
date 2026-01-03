package service

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/teacinema-go/gateway-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Manager struct {
	Auth     AuthService
	services []io.Closer
}

func NewServiceManager(serviceConfig config.ServiceConfig, logger *slog.Logger) (*Manager, error) {
	authService, err := NewAuthService(serviceConfig.AuthServiceURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	return &Manager{
		Auth:     authService,
		services: []io.Closer{authService},
	}, nil
}

func (m *Manager) Close() error {
	var firstErr error
	for _, service := range m.services {
		if err := service.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func newGRPCConnection(serviceURL string) (*grpc.ClientConn, error) {
	// TODO: use TLS in production
	conn, err := grpc.NewClient(
		serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect at %s: %w", serviceURL, err)
	}

	return conn, nil
}
