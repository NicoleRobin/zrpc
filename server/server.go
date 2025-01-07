package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"github.com/nicolerobin/zrpc/log"
)

// Serve start service
func Serve() error {
	ctx := context.Background()
	s := newServer()

	return start(ctx, s)
}

func start(ctx context.Context, s *grpc.Server) error {
	l, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	return s.Serve(l)
}
