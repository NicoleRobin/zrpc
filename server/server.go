package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/nicolerobin/zrpc/log"
	"google.golang.org/grpc"
)

type serverWrapper struct {
	s *grpc.Server
}

func (s *serverWrapper) Serve(l net.Listener) error {
	return s.serveGrpc(l)
}

func (s *serverWrapper) register(service serviceInfo) {
	service.r.RegisterService(s.s, service.handler)
}

func (s *serverWrapper) serveGrpc(l net.Listener) error {
	return s.s.Serve(l)
}

func (s *serverWrapper) serveHttp(ctx context.Context) error {
	l, err := net.Listen("tcp", "")
	if err != nil {
		return fmt.Errorf("net.Listen() failed: %w", err)
	}
	return http.Serve(l, s.s)
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
		s: NewGrpcServer(Option{GrpcOptions: options}),
	}
}

// Start run server
func Start(ctx context.Context) error {
	l, err := net.Listen("tcp", "")
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof(ctx, "starting server on %s", l.Addr())
	s := newServer()
	registerServices(s)
	return s.Serve(l)
}
