package service

import (
	"context"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/dto"
	"google.golang.org/grpc"
)

type authService struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthService(authServiceUrl string) (AuthService, error) {
	conn, err := newGRPCConnection(authServiceUrl)
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)
	return &authService{
		client: client,
		conn:   conn,
	}, nil
}

func (s *authService) Close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) SendOtp(ctx context.Context, r dto.SendOtpRequest) (*authv1.SendOtpResponse, error) {
	req := &authv1.SendOtpRequest{
		Identifier: r.Identifier,
		Type:       string(r.Type),
	}

	res, err := s.client.SendOtp(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
