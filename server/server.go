package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Serve start service
func Serve() error {
	s := newServer()

	return start(s)
}

func start(s *grpc.Server) error {
	l, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return s.Serve(l)
}
