package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type rpcServiceRegister interface {
	RegisterService(s grpc.ServiceRegistrar, handler interface{})
}

func RegisterService(r rpcServiceRegister, handler interface{}) {
}

func handlerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return nil, nil
}

type Option struct {
	GrpcOptions []grpc.ServerOption
}

type serverWrapper struct {
	s *grpc.Server
}

func (s *serverWrapper) Serve(l net.Listener) error {
	return s.serveGrpc(l)
}

func (s *serverWrapper) serveGrpc(l net.Listener) error {
	return s.s.Serve(l)
}

func (s *serverWrapper) serveHttp(ctx context.Context) error {
	return nil
}

func newServer(options ...grpc.ServerOption) *serverWrapper {
	return &serverWrapper{
		s: NewServer(Option{GrpcOptions: options}),
	}
}

func NewServer(opt Option) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(handlerInterceptor),
	}
	opt.GrpcOptions = append(opts, opt.GrpcOptions...)

	return nil
}
