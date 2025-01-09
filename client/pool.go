package client

import (
	"context"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
)

type ccEntry struct {
	r  atomic.Value
	cc atomic.Value

	done     chan struct{}
	creating uint32

	wg       sync.WaitGroup
	resolver atomic.Value
}

func (e *ccEntry) assign(cc *grpc.ClientConn) {
}

func (e *ccEntry) pick() *grpc.ClientConn {
	cc, _ := e.cc.Load().(*grpc.ClientConn)
	return cc
}

var (
	ccMap sync.Map
)

func getCc(ctx context.Context, name string) (*grpc.ClientConn, bool) {
	v, ok := ccMap.Load(name)
	if !ok {
		return nil, false
	}

	cc := v.(*ccEntry).pick()
	if cc != nil {
		return cc, true
	}

	return nil, false
}
