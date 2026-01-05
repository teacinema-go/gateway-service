package clients

import (
	"context"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/dto"
	"google.golang.org/grpc"
)

type authServiceClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

type AuthServiceClient interface {
	SendOtp(ctx context.Context, req dto.SendOtpRequest) (*authv1.SendOtpResponse, error)
	Close() error
}

func NewAuthServiceClient(authServiceUrl string) (AuthServiceClient, error) {
	conn, err := newGRPCConnection(authServiceUrl)
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)
	return &authServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (s *authServiceClient) Close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *authServiceClient) SendOtp(ctx context.Context, r dto.SendOtpRequest) (*authv1.SendOtpResponse, error) {
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
