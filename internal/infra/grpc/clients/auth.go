package clients

import (
	"context"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/auth/dto/request"
	"google.golang.org/grpc"
)

type authServiceClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
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

func (s *authServiceClient) SendOtp(ctx context.Context, r request.SendOtpRequest) (*authv1.SendOtpResponse, error) {
	req := &authv1.SendOtpRequest{
		Identifier:     r.Identifier,
		IdentifierType: r.IdentifierType.ToProto(),
	}

	res, err := s.client.SendOtp(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *authServiceClient) VerifyOtp(ctx context.Context, r request.VerifyOtpRequest) (*authv1.VerifyOtpResponse, error) {
	req := &authv1.VerifyOtpRequest{
		Identifier:     r.Identifier,
		IdentifierType: r.IdentifierType.ToProto(),
		Otp:            r.Otp,
	}

	res, err := s.client.VerifyOtp(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *authServiceClient) Refresh(ctx context.Context, refreshToken string) (*authv1.RefreshResponse, error) {
	req := &authv1.RefreshRequest{
		RefreshToken: refreshToken,
	}

	res, err := s.client.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *authServiceClient) Logout(ctx context.Context, refreshToken string) (*authv1.LogoutResponse, error) {
	req := &authv1.LogoutRequest{
		RefreshToken: refreshToken,
	}

	res, err := s.client.Logout(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
