package clients

import (
	"fmt"

	"github.com/teacinema-go/gateway-service/internal/config"
)

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
