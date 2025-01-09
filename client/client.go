package client

import (
	"context"
	"github.com/nicolerobin/zrpc/config"
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
	return m.GetClientWithConfig(ctx, config.ClientConfig{})
}

func (m *Manager) GetClientWithConfig(ctx context.Context, conf config.ClientConfig) (*grpc.ClientConn, error) {
	cc, ok := getCc(m.clientName)
	if ok {
		return cc, nil
	}

	cc, err := m.dail(ctx, conf)
	if err == nil {
		return cc, nil
	}

	return nil, err
}

func (m *Manager) dail(ctx context.Context, conf config.ClientConfig) (*grpc.ClientConn, error) {
	return nil, nil
}

func (m *Manager) Close(error) {}

func Get(clientName string) *Manager {
	return &Manager{
		clientName:      clientName,
		dialBreakerName: "rpc-client-dial-" + clientName,
	}
}
