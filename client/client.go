package client

import (
	"context"
	"fmt"
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
	cc, ok := getCc(ctx, m.clientName)
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
	return createCc(ctx, m.clientName)
}

func (m *Manager) dailWithBreaker(ctx context.Context, conf config.ClientConfig, cb func(cc *grpc.ClientConn) error) error {
	return nil
}

func (m *Manager) Close(error) {}

func Get(clientName string) *Manager {
	return &Manager{
		clientName:      clientName,
		dialBreakerName: "rpc-client-dial-" + clientName,
	}
}

func createCc(ctx context.Context, name string) (*grpc.ClientConn, error) {
	dsn := fmt.Sprintf("dns:///%s", name)

	grpcClient, err := grpc.NewClient(dsn)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient() failed, err: %w", err)
	}

	return grpcClient, nil
}
