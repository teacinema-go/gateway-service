package grpc

import (
	"context"
	"log/slog"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/gateway-service/internal/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
	l      *slog.Logger
}

func NewAuthClient(authServiceUrl string, l *slog.Logger) (*AuthClient, error) {
	// TODO use TLS in prod
	conn, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)
	return &AuthClient{
		client: client,
		conn:   conn,
		l:      l,
	}, nil
}

func (c *AuthClient) Close() error {
	c.l.Info("Closing Auth Client")
	err := c.conn.Close()
	if err != nil {
		c.l.Error("Failed to close Auth Client")
		return err
	}
	c.l.Info("Auth Client closed")
	return nil
}

func (c *AuthClient) SendOtp(ctx context.Context, r dto.SendOtpRequest) (*authv1.SendOtpResponse, error) {
	req := &authv1.SendOtpRequest{
		Identifier: r.Identifier,
		Type:       string(r.Type),
	}
	c.l.Info("Sending Otp request...")
	res, err := c.client.SendOtp(ctx, req)
	if err != nil {
		c.l.Error("Failed to send Otp request")
		return nil, err
	}
	c.l.Info("Otp sent successfully")
	return res, nil
}
