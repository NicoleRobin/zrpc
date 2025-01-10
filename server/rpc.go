package server

import (
	"context"
	"github.com/nicolerobin/zrpc/log"
	"google.golang.org/grpc"
	"sync"
)

var (
	servicesLock sync.Mutex

	grpcServices []serviceInfo
)

type rpcServiceRegister interface {
	RegisterService(s grpc.ServiceRegistrar, handler interface{})
}

type serviceInfo struct {
	r       rpcServiceRegister
	handler interface{}
}

func RegisterService(r rpcServiceRegister, handler interface{}) {
	servicesLock.Lock()
	defer servicesLock.Unlock()

	grpcServices = append(grpcServices, serviceInfo{
		r:       r,
		handler: handler,
	})
}

func handlerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Info(ctx, "entrance")
	return handler(ctx, req)
}

type Option struct {
	GrpcOptions []grpc.ServerOption
}

func NewGrpcServer(opt Option) *grpc.Server {
	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(handlerInterceptor),
	}
	options = append(options, opt.GrpcOptions...)

	server := grpc.NewServer(options...)

	return server
}
