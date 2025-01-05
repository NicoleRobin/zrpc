package client

import (
	"context"
	"google.golang.org/grpc"
)

type Client struct {
}

func (c *Client) Invoke(ctx context.Context, path string, in, out interface{}, opts ...grpc.CallOption) error {
	return nil
}

type Manager struct {
	clientName      string
	dialBreakerName string
}

func (m *Manager) GetClient(ctx context.Context) (*grpc.ClientConn, error) {
	return nil, nil
}

func (m *Manager) Close(error) {}

func Get(clientName string) *Manager {
	return &Manager{
		clientName:      clientName,
		dialBreakerName: "rpc-client-dial-" + clientName,
	}
}
