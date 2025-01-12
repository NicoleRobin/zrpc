package server

import (
	"context"
	"fmt"
	"github.com/nicolerobin/zrpc/config"
	"github.com/nicolerobin/zrpc/core"
	"github.com/nicolerobin/zrpc/log"
	"google.golang.org/grpc"
	"net"
)

type serverWrapper struct {
	grpcServer *grpc.Server
}

func (s *serverWrapper) Serve(ctx context.Context, l net.Listener) error {
	return s.serveGrpc(ctx, l)
}

func (s *serverWrapper) register(service serviceInfo) {
	service.r.RegisterService(s.grpcServer, service.handler)
}

func (s *serverWrapper) serveGrpc(ctx context.Context, l net.Listener) error {
	return s.grpcServer.Serve(l)
}

func registerServices(s *serverWrapper) {
	servicesLock.Lock()
	for _, service := range grpcServices {
		s.register(service)
	}
	servicesLock.Unlock()
}

func newServer(options ...grpc.ServerOption) *serverWrapper {
	return &serverWrapper{
		grpcServer: NewGrpcServer(Option{GrpcOptions: options}),
	}
}

// Start run server
func Start(ctx context.Context) error {
	l, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof(ctx, "starting server on %s", l.Addr())
	s := newServer()
	registerServices(s)
	return s.Serve(ctx, l)
}

func beforeCall(ctx context.Context, req interface{}) (context.Context, error) {
	return ctx, nil
}

func init() {
	core.AddHook(Hook, BeforeCall(beforeCall))
}
