package service

import (
	"context"
	"log/slog"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/dto"
	"google.golang.org/grpc"
)

type authService struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
	l      *slog.Logger
}

func NewAuthService(authServiceUrl string, l *slog.Logger) (AuthService, error) {
	conn, err := newGRPCConnection(authServiceUrl)
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)
	return &authService{
		client: client,
		conn:   conn,
		l:      l,
	}, nil
}

func (s *authService) Close() error {
	s.l.Info("Closing Auth Client")
	err := s.conn.Close()
	if err != nil {
		s.l.Error("Failed to close Auth Client")
		return err
	}
	s.l.Info("Auth Client closed")
	return nil
}

func (s *authService) SendOtp(ctx context.Context, r dto.SendOtpRequest) (*authv1.SendOtpResponse, error) {
	req := &authv1.SendOtpRequest{
		Identifier: r.Identifier,
		Type:       string(r.Type),
	}
	s.l.Info("Sending Otp request")
	res, err := s.client.SendOtp(ctx, req)
	if err != nil {
		s.l.Error("Failed to send Otp request")
		return nil, err
	}
	s.l.Info("Otp sent successfully")
	return res, nil
}
