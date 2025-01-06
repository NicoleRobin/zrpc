package server

import (
	"context"

	"google.golang.org/grpc"
)

// Serve start service
func Serve() error {
	s := newServer()

	return Start(s)
}

func Start(s *grpc.Server) error {
	ctx := context.Background()

	s.Serve()

	return nil
}
