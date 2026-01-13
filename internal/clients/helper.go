package clients

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGRPCConnection(serviceURL string) (*grpc.ClientConn, error) {
	// TODO: use TLS in production, use circuit breaker or other pattern in production
	conn, err := grpc.NewClient(
		serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect at %s: %w", serviceURL, err)
	}

	return conn, nil
}
