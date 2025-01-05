package rpc

import (
	"context"

	"google.golang.org/grpc"
)

type ClientGetter interface {
	GetClient(context.Context) (*grpc.ClientConn, error)
	Close(error)
}

func NewRawGetter(cc grpc.ClientConnInterface) ClientGetter {
	return &rawGetter{cc: cc}
}

type rawGetter struct {
	cc grpc.ClientConnInterface
}

func (r *rawGetter) GetClient(ctx context.Context) (*grpc.ClientConn, error) {
	return r.cc.(*grpc.ClientConn), nil
}

func (r *rawGetter) Close(err error) {
}
