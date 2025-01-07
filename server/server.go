package server

import (
	"context"
	"fmt"
	"github.com/nicolerobin/zrpc/config"
	"net"

	"github.com/nicolerobin/zrpc/log"
)

// Start run server
func Start(ctx context.Context) error {
	l, err := net.Listen("tcp", config.GetAddr())
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	s := newServer()
	return s.Serve(l)
}
